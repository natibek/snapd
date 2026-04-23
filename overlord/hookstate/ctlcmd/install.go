// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2024 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package ctlcmd

import (
	"fmt"

	"github.com/snapcore/snapd/i18n"
	"github.com/snapcore/snapd/snap"
	"github.com/snapcore/snapd/strutil"
)

var (
	shortInstallHelp = i18n.G("Install components")
	longInstallHelp  = i18n.G(`
The install command installs components.
`)
)

func init() {
	addCommand("install", shortInstallHelp, longInstallHelp, func() command { return &installCommand{} })
}

type installCommand struct {
	baseCommand
	Positional struct {
		Names []string `positional-arg-name:"<snap|snap+comp|+comp>" required:"yes" description:"Components to be installed (snap must be the caller snap if specified)."`
	} `positional-args:"yes"`
}

func (c *installCommand) Execute([]string) error {
	ctx, err := c.ensureContext()
	if err != nil {
		return err
	}

	comps, err := validateSnapAndCompsNames(c.Positional.Names, ctx.InstanceName())
	if err != nil {
		return err
	}

	affected, err := runSnapManagementCommand(ctx, managementCommand{
		operation: installManagementCommand, components: comps})

	if err != nil {
		_, ok := err.(*snap.AlreadyInstalledError)
		if !ok {
			return err
		}
	}

	if len(affected) < len(comps) {
		for _, comp := range comps {
			if !strutil.ListContains(affected, comp) {
				msg := fmt.Sprintf(i18n.G(`snapctl: component %q is already installed`), comp)
				fmt.Fprintln(c.stderr, msg)
			}
		}
	}

	return nil
}
