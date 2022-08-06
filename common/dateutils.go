package common

import "time"

func CurStdDate() time.Time {
	loc,_:=time.LoadLocation(StdLocation)
	t:=time.Now()
	t = t.In(loc)
	return t
}

func CurLocalDate() time.Time{
	loc,_:=time.LoadLocation(LocalLocation)
	t:=time.Now()
	t = t.In(loc)
	return t
}

func CurCstDate() time.Time{
	t:=time.Now()
	t = t.In(CSTZone)
	return t
}