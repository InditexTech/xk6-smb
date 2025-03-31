// SPDX-FileCopyrightText: 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

package xk6smb

import (
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/smb", new(Client))
}

var logger = logrus.New()
