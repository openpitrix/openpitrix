// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func getJwtCmd() *cobra.Command {
	var (
		userId     string
		role       string
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
			if len(role) == 0 {
				return fmt.Errorf("[role] should specify")
			}
			if expireTime == 0 {
				return fmt.Errorf("[expire_time] should specify")
			}
			if len(secretKey) == 0 {
				return fmt.Errorf("[secret_key] should specify")
			}

			token, err := senderutil.Generate(secretKey, expireTime, userId, role)
			if err != nil {
				return err
			}
			fmt.Print(token)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&userId, "user_id", "u", "", "specify user_id in JWT")
	f.StringVarP(&role, "role", "r", "", "specify role in JWT")
	f.DurationVarP(&expireTime, "expire_time", "e", 2*time.Hour, "specify expire_time in JWT")
	f.StringVarP(&secretKey, "secret_key", "s", "", "specify secret_key in JWT")

	return cmd
}
