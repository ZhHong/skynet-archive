//Copyright (c) 2011 Brian Ketelsen

//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package skylib

import (
	"fmt"
	"net/rpc"
	"time"
)

// RpcService is a struct that represents a remotly 
// callable function.  It is intented to be part of 
// an array or collection of RpcServices.  It contains
// a member "Provides" which is the name of the service the
// remote call provides, and a Client pointer which is a pointer
// to an RPC client connected to this service.
type RpcService struct {
	Provides string
}

func (r *RpcService) parseError(err string) {
	panic(&Error{err, r.Provides})
}

// A Generic struct to represent any service in the SkyNet system.
type Service struct {
	IPAddress string
	Name      string
	Port      int
	Provides  string
}

// A HeartbeatRequest is the struct that is sent for ping checks.
type HeartbeatRequest struct {
	Timestamp int64
}

// HeartbeatResponse is the struct that is returned on a ping check.
type HeartbeatResponse struct {
	Timestamp time.Time
	Ok        bool
}

// HealthCheckRequest is the struct that is sent on a more advanced heartbeat request.
type HealthCheckRequest struct {
	Timestamp time.Time
}

// HealthCheckResponse is the struct that is sent back to the HealthCheckRequest-er
type HealthCheckResponse struct {
	Timestamp time.Time
	Load      float64
}

// The struct that is stored in the Route
// Async delineates whether it's ok to call this and not
// care about the response.
// OkToRetry delineates whether it's ok to call this service
// more than once.
type RpcCall struct {
	Service   string
	Async     bool
	OkToRetry bool
	ErrOnFail bool
}

// Parent struct for the configuration
type NetworkServers struct {
	Services []*Service
}

type ServerConfig interface {
	Equal(that interface{}) bool
}

// Exported RPC method for the health check
func (hc *Service) Ping(hr *HeartbeatRequest, resp *HeartbeatResponse) (err error) {

	resp.Timestamp = time.Now()

	return nil
}

// Exported RPC method for the advanced health check
func (hc *Service) PingAdvanced(hr *HealthCheckRequest, resp *HealthCheckResponse) (err error) {

	resp.Timestamp = time.Now()
	resp.Load = 0.1 //todo
	return nil
}

// Method to register the heartbeat of each skynet
// client with the healthcheck exporter.
func RegisterHeartbeat() {
	r := NewService("Service.Ping")
	rpc.Register(r)
}

func (r *Service) Equal(that *Service) bool {
	var b bool
	b = false
	if r.Name != that.Name {
		return b
	}
	if r.IPAddress != that.IPAddress {
		return b
	}
	if r.Port != that.Port {
		return b
	}
	if r.Provides != that.Provides {
		return b
	}
	b = true
	return b
}

// Utility function to return a new Service struct
// pre-populated with the data on the command line.
func NewService(provides string) *Service {
	r := &Service{
		Name:      *Name,
		Port:      *Port,
		IPAddress: *BindIP,
		Provides:  provides,
	}

	return r
}

type Error struct {
	Msg     string
	Service string
}

func (e *Error) Error() string { return fmt.Sprintf("Service %s had error: %s", e.Service, e.Msg) }

func NewError(msg string, service string) (err *Error) {
	err = &Error{Msg: msg, Service: service}
	return
}
