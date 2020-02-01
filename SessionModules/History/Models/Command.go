package History

/*
I know this struct is repeated but that is because later on
I want to make a plug-in system so that people create their own plug ins
and each plug-in will create their service and models and models will repeat

Also this commands could be different from the other ones
*/
type Command struct {
	Cmnd string
}

func NewCommand(Cmnd string) (c Command) {
	c.Cmnd = Cmnd
	return
}
