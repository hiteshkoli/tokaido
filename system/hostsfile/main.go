package hostsfile

import (
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
)

const localhost = "127.0.0.1"

func confirmRemove(hostname string) bool {
	c := utils.ConfirmationPrompt("Would you like Tokaido to automatically update your hosts file, removing the host '"+hostname+"'? You may be prompted for elevated access.", "y")
	if c == false {
		fmt.Println(`Your hostsfile can be amended manually in order to remove this hostname. See https://tokaido.io/docs/config/#updating-your-hostsfile for more information.`)
	}

	return c
}

// RemoveEntry - Remove an entry from /etc/hosts or equivalent
func RemoveEntry(hostname string) error {
	hosts, err := NewHosts()
	if err != nil {
		return err
	}

	if hosts.Has(localhost, hostname) {
		if confirmRemove(hostname) == false {
			return nil
		}

		hosts.Remove(localhost, hostname)
		if hosts.IsWritable() == false {
			err := hosts.WriteElevated()
			if err != nil {
				return err
			}

			return nil
		}

		if err := hosts.Flush(); err != nil {
			return err
		}

		return nil
	}

	return nil
}
