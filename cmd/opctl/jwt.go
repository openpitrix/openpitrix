// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"openpitrix.io/openpitrix/pkg/util/jwtutil"
)

func getJwtCmd() *cobra.Command {
	var (
		userId     string
		expireTime time.Duration
		secretKey  string
	)
	cmd := &cobra.Command{
		Use:   "generate_jwt",
		Short: "Generate the JWT",
		Long:  "Generate the JWT with specify secret key and user_id and role and expire_time",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Flags().Parse(args)
			if err != nil {
				return err
			}

			if len(userId) == 0 {
				return fmt.Errorf("[user_id] should specify")
			}
			if expireTime == 0 {
				return fmt.Errorf("[expire_time] should specify")
			}
			if len(secretKey) == 0 {
				return fmt.Errorf("[secret_key] should specify")
			}

			token, err := jwtutil.Generate(secretKey, expireTime, userId)
			if err != nil {
				return err
			}
			fmt.Print(token)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&userId, "user_id", "u", "", "specify user_id in JWT")
	f.DurationVarP(&expireTime, "expire_time", "e", 2*time.Hour, "specify expire_time in JWT")
	f.StringVarP(&secretKey, "secret_key", "s", "", "specify secret_key in JWT")

	return cmd
}
