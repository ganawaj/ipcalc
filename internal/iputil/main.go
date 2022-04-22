package iputil

import (
	"fmt"

	"inet.af/netaddr"
)

func GetIPPrefix(addr string, mask string) (netaddr.IPPrefix, error){

	return netaddr.ParseIPPrefix(fmt.Sprintf("%s/%s", addr, mask))

}

func GetBroadcastAddr (addr string, mask string) (string, error){

	prefix, err := netaddr.ParseIPPrefix(fmt.Sprintf("%s/%s", addr, mask))
	if err != nil {
		return "", err
	}

	return prefix.Range().To().String(), nil

}

func GetHostMin (addr string, mask string) (string, error){

	prefix, err := netaddr.ParseIPPrefix(fmt.Sprintf("%s/%s", addr, mask))
	if err != nil {
		return "", err
	}

	return prefix.Range().From().Next().String(), nil
}

func GetHostMax  (addr string, mask string) (string, error){

	prefix, err := netaddr.ParseIPPrefix(fmt.Sprintf("%s/%s", addr, mask))
	if err != nil {
		return "", err
	}

	return prefix.Range().To().Prior().String(), nil
}

// func (service) GetNetworkInfo (addr string, mask string) (IP, error){

// 	if addr == "" || mask == "" {
// 		return IP{}, ErrEmpty
// 	}

// 	// create prefix from ip address and range
// 	prefix, err := netaddr.ParseIPPrefix(fmt.Sprintf("%s/%s", addr, mask))
// 	if err != nil {
// 		return IP{}, ErrEmpty
// 	}

// 	return IP{
// 		Address: 	addr,
// 		Netmask: 	mask,
// 		Network: 	fmt.Sprintf("%s/%s", prefix.Range().From().String(), mask),
// 		Broadcast: 	prefix.Range().To().String(),
// 		HostMin: 	prefix.Range().From().Next().String(),
// 		HostMax: 	prefix.Range().To().Prior().String(),
// 	}, nil

// }