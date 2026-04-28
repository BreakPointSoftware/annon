package annon

func (a *Anonymiser) Copy(input any) (any, error) { return a.walker.Copy(input) }
