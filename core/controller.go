package core

type Controller struct {
	writer  *Writer
	reader  *Reader
	signals chan int
}

func NewController(writer *Writer, reader *Reader) *Controller {
	controller := &Controller{
		writer:  writer,
		reader:  reader,
		signals: make(chan int, 100),
	}
	return controller
}

func (c *Controller) Run() {

}
