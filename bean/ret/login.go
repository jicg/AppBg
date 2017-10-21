package ret

func Login_Error(msg string) *Ret {
	return Error(msg)
}

func Login_Success(token string) *Ret {
	return Success(H{"token": token})
}
