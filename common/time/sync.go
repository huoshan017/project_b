package time

const (
	NetworkDelayListLength = 100
)

type delayList struct {
	list        []Duration
	writeOffset int
	readOffset  int
	write       int
}

func createDelayList(l int) *delayList {
	return &delayList{
		list: make([]Duration, l),
	}
}

func (dl *delayList) push(delay Duration) {
	// 满了就丢掉最早的一个
	if dl.write >= len(dl.list) {
		dl.pop()
	}
	dl.list[dl.writeOffset] = delay
	dl.writeOffset += 1
	if dl.writeOffset >= len(dl.list) {
		dl.writeOffset = 0
	}
	dl.write += 1
}

func (dl *delayList) pop() {
	if dl.write <= 0 {
		return
	}
	dl.readOffset += 1
	if dl.readOffset >= len(dl.list) {
		dl.readOffset = 0
	}
	dl.write -= 1
}

func (dl *delayList) getDelay() (Duration, bool) {
	if dl.write == 0 {
		return 0, false
	}
	var duration Duration
	if dl.writeOffset == 0 {
		duration = dl.list[len(dl.list)-1]
	} else {
		duration = dl.list[dl.writeOffset-1]
	}
	return duration, true
}

func (dl *delayList) getAvgDelay() (Duration, bool) {
	if dl.write == 0 {
		return 0, false
	}
	var total Duration
	if dl.writeOffset > dl.readOffset {
		for i := dl.readOffset; i < dl.writeOffset; i++ {
			total += dl.list[i]
		}
	} else if dl.writeOffset < dl.readOffset {
		for i := dl.readOffset; i < len(dl.list); i++ {
			total += dl.list[i]
		}
		for i := 0; i < dl.writeOffset; i++ {
			total += dl.list[i]
		}
	} else {
		for i := 0; i < len(dl.list); i++ {
			total += dl.list[i]
		}
	}
	return total / Duration(dl.write), true
}

type Client struct {
	serverTime         CustomTime
	sendTime, recvTime CustomTime
	waitRecv           bool
	delayList          *delayList
}

func NewClient() *Client {
	return &Client{
		delayList: createDelayList(NetworkDelayListLength),
	}
}

func (c *Client) GetSendTime() CustomTime {
	return c.sendTime
}

func (c *Client) SetSendTime(t CustomTime) {
	if c.waitRecv {
		return
	}
	c.sendTime = t
	c.waitRecv = true
}

func (c *Client) SetRecvTime(t CustomTime) {
	if !c.waitRecv {
		return
	}
	c.recvTime = t
	c.waitRecv = false
}

// 如果设置服务器时间，则要跟接收时间一起设置，确保同时变化
func (c *Client) SetRecvAndServerTime(rt, st CustomTime) {
	c.SetRecvTime(rt)
	c.serverTime = st
	c.delayList.push(rt.Sub(c.sendTime) / 2)
}

// 客户端与服务器的时间差
func (c *Client) GetDeltaC2S() (Duration, bool) {
	if c.waitRecv {
		return 0, false
	}

	return c.recvTime.Sub(c.serverTime), true
}

func (c *Client) GetDelay() Duration {
	d, r := c.delayList.getDelay()
	if !r {
		return -1
	}
	return d
}

func (c *Client) GetAvgDelay() Duration {
	d, r := c.delayList.getAvgDelay()
	if !r {
		return -1
	}
	return d
}
