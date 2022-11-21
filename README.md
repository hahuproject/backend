# hahu_sms_backend

A backend for hahu's school management system

pkg -> ports | adapters
ports -> core
package ports
type AuthPort interface{
Register()
Login()
}
adapters -> core -> auth -> auth.go
package auth
type Adapter {}
func NewAdapter() \*Adapter {return &Adapter{}}
//implement AuthPort interface
