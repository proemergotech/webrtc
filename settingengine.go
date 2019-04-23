// +build !js

package webrtc

import (
	"net"
	"time"

	"github.com/pion/ice"
	"github.com/pion/logging"
)

// SettingEngine allows influencing behavior in ways that are not
// supported by the WebRTC API. This allows us to support additional
// use-cases without deviating from the WebRTC API elsewhere.
type SettingEngine struct {
	ephemeralUDP struct {
		PortMin uint16
		PortMax uint16
	}
	detach struct {
		DataChannels bool
	}
	timeout struct {
		ICEConnection *time.Duration
		ICEKeepalive  *time.Duration
	}
	candidates struct {
		ICENetworkTypes []NetworkType
		localNatRule    NatRule
	}
	LoggerFactory logging.LoggerFactory
}

type NatRule func(network string, localIP net.IP, localPort int) (natIP net.IP, natPort int)

// DetachDataChannels enables detaching data channels. When enabled
// data channels have to be detached in the OnOpen callback using the
// DataChannel.Detach method.
func (e *SettingEngine) DetachDataChannels() {
	e.detach.DataChannels = true
}

// SetConnectionTimeout sets the amount of silence needed on a given candidate pair
// before the ICE agent considers the pair timed out.
func (e *SettingEngine) SetConnectionTimeout(connectionTimeout, keepAlive time.Duration) {
	e.timeout.ICEConnection = &connectionTimeout
	e.timeout.ICEKeepalive = &keepAlive
}

// SetEphemeralUDPPortRange limits the pool of ephemeral ports that
// ICE UDP connections can allocate from. This setting currently only
// affects host candidates, not server reflexive candidates.
func (e *SettingEngine) SetEphemeralUDPPortRange(portMin, portMax uint16) error {
	if portMax < portMin {
		return ice.ErrPort
	}

	e.ephemeralUDP.PortMin = portMin
	e.ephemeralUDP.PortMax = portMax
	return nil
}

// SetNetworkTypes configures what types of candidate networks are supported
// during local and server reflexive gathering.
func (e *SettingEngine) SetNetworkTypes(candidateTypes []NetworkType) {
	e.candidates.ICENetworkTypes = candidateTypes
}

// SetLocalNatRule lets you publish different ip/ports as local ip than the ones you listen to. Useful for lan NAT,
// eg: when running the webrtc server inside a container, but joining to it using
func (e *SettingEngine) SetLocalNatRule(rule NatRule) {
	e.candidates.localNatRule = rule
}
