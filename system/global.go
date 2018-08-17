package system

var UserCache *MemCache
var LockManager = newLockerManager()
const TokenPath  = "/fxgos/token"
const ProtectedPrefix  = "/fxgos"
const PublicPrefix  = "/public"

var SharedLocker  = &SharedLock{
	Module: "module",
}
