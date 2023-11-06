package secrets

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/xopoww/standup/internal/auth"
	"github.com/xopoww/standup/internal/standupctl/util"
)

func genKeys() *cobra.Command {
	var args struct {
		dir     string
		crtFile string // cert == public key
		keyFile string // key == private key
		force   bool
	}
	cmd := &cobra.Command{
		Use:     "generate-keys",
		Aliases: []string{"gen-keys"},
		Short:   "Generate key pair",
		Long:    "Generate new private-public key pair (for jwt)",
		RunE: func(_ *cobra.Command, _ []string) error {
			if args.dir == "" {
				cwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("cwd: %w", err)
				}
				args.dir = cwd
			}
			args.crtFile = absolutePath(args.crtFile, args.dir, "--crt")
			args.keyFile = absolutePath(args.keyFile, args.dir, "--key")
			return runGenKeys(args.crtFile, args.keyFile, args.force)
		},
	}
	cmd.Flags().StringVar(&args.dir, "dir", "",
		"output directory (default is current working directory)",
	)
	cmd.Flags().StringVar(&args.crtFile, "crt", "public.pem",
		"public key (certificate) filename (absolute path overrides '-dir')",
	)
	cmd.Flags().StringVar(&args.keyFile, "key", "private.ec.key",
		"private key filename (absolute path overrides '-dir')",
	)
	cmd.Flags().BoolVarP(&args.force, "force", "f", false, "overwrite existing files without confirmation")

	return cmd
}

func runGenKeys(crtFile, keyFile string, force bool) error {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("generate key: %w", err)
	}

	kf, err := createOrTruncate(keyFile, force)
	if errors.Is(err, util.ErrAborted) {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return nil
	} else if err != nil {
		return fmt.Errorf("open key file: %w", err)
	}
	err = auth.WritePrivateKey(key, kf)
	if err != nil {
		return fmt.Errorf("write key: %w", err)
	}
	err = kf.Close()
	if err != nil {
		return fmt.Errorf("save key: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Saved private key to %q.\n", keyFile)

	cf, err := createOrTruncate(crtFile, force)
	if errors.Is(err, util.ErrAborted) {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		return nil
	} else if err != nil {
		return fmt.Errorf("open crt file: %w", err)
	}
	err = auth.WritePublicKey(&key.PublicKey, cf)
	if err != nil {
		return fmt.Errorf("write crt: %w", err)
	}
	err = cf.Close()
	if err != nil {
		return fmt.Errorf("save crt: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Saved public key to %q.\n", crtFile)

	return nil
}

func absolutePath(file, dir, flag string) string {
	if path.IsAbs(file) {
		fmt.Fprintf(os.Stderr, "Warning: '%s' is abosulte (%q) which overrides '--dir'.\n", flag, file)
		return file
	}
	return path.Join(dir, file)
}

func createOrTruncate(file string, force bool) (*os.File, error) {
	finfo, err := os.Stat(file)
	//nolint:nestif // seems ok to me ¯\_(ツ)_/¯
	if err == nil {
		if finfo.IsDir() {
			return nil, fmt.Errorf("%q is a directory", file)
		}
		if !force {
			err := util.Confirmf(os.Stderr, os.Stdin, "Overwrite %q?", file)
			if err != nil {
				return nil, err
			}
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("stat: %w", err)
	}
	return os.Create(file)
}
