package ret

func Token_Refresh(msg string, token string) *Ret {
	return &Ret{
		Code: -1001,
		Msg:  msg,
		Data: H{"token": token},
	}
}

func Token_Error(msg string) *Ret {
	return &Ret{
		Code: -1,
		Msg:  msg,
	}
}
