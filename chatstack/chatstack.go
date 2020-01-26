package chatstack

type ChatMessage struct {
	From    string `json:'from'`
	To      string `json:'to',omitempty` // empty = world/region
	Message string `json:'msg'`
}

var MaxChatSize = 255

type (
	ChatStack struct {
		top    *chatnode
		length int
	}
	chatnode struct {
		value ChatMessage
		prev  *chatnode
	}
)

// Create a new ChatStack
func New() *ChatStack {
	return &ChatStack{nil, 0}
}

// Return the number of items in the ChatStack
func (this *ChatStack) Len() int {
	return this.length
}

// View the top item on the ChatStack
func (this *ChatStack) Peek() ChatMessage {
	if this.length == 0 {
		return ChatMessage{}
	}
	return this.top.value
}

// Pop the top item of the ChatStack and return it
func (this *ChatStack) Pop() ChatMessage {
	if this.length == 0 {
		return ChatMessage{}
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

// Push a value onto the top of the ChatStack
func (this *ChatStack) Push(value ChatMessage) {
	if this.length > MaxChatSize {
		println("Chat buffer full, waiting")
		this.Trim()
	}
	n := &chatnode{value, this.top}
	this.top = n
	this.length++
}

func (this *ChatStack) Trim() {
	var n, m *chatnode
	n = this.top

	for i := 0; i < this.length; i++ {
		n = n.prev
		m = n
		if i > MaxChatSize {
			if m == nil {
				break
			}
			m.prev = nil
			this.length--
		}
	}

}
