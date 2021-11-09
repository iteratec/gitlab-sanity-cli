# The AbstractHandler

## Interface

The AbstractHandler interface require three functions to be implemented

```go
type AbstractHandler interface {
	Pipe(*sync.WaitGroup, chan<- content, <-chan int)
	Controller(ChannelHandlerFunc)
	ApiCallFunc(int) ([]int, *gitlab.Response, error)
}
```

## The Pipe() function

The `Pipe()` function will instantiate the goroutines from inside `genericHandler()` function and requires three parameters.
- `waitgroup` pointer is used to mark goroutine itself as done when the Pipe() function ends.
- `content` channel is of type struct and is used as channel for the `ChannelHandlerFunc` (egress channel)
- `int` channel is used to get the current gitlab API call result page number. The page number is send from `genericHandler()` function.

To use the Pipe() function, you must interprete the `int` channel and put results into the `content` channel.

Example:

```go
func (h MyHandler) Pipe(wg *sync.WaitGroup, output chan<- content, page <-chan int))
{
	defer wg.Done()
	for id := range page {
		res, _, err := gitlab.Get....(id)
		if err != nil {
			log.Fatal(err)
		}
		output <- content{
			Req:  "MyReq",
			ID:   res.ID,
		}
	}
}
```

## The Controller() function

The `Controller()` function has only one Parameter: `ChannelHandlerFunc` and is called from function `Run()` in AbstractHandler.
- Parameter `ChannelHandlerFunc` is used to handle the results from the `content` channel.


The `Controller()` function can be used to implement your own handler (like `genericHandler()`).
In most cases, you just call `genericHandler()` at the end of `Controller()`.
If you don't want to call `genericHandler()` from `Controller()`, you have to implement your own "genericHandler" function which calls the `ApiCallFunc()` and `Pipe()`

Example:

```go
func (h MyHandler) Controller(channelHandlerFunc ChannelHandlerFunc) {
	_, resp, err := gitlab.Users.ListUsers(&gitlab.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	genericHandler(h, channelHandlerFunc, resp.TotalPages)
}
```

## The ApiCallFunc() function

The `ApiCallFunc()` function has only one Parameter and returns three types.
- Parameter: `int` is used to get the current gitlab API call result Page number. The Page number is send from `genericHandler()` function.
- Return `[]int` slice is used as ID array and needs to contains the gitlab API Call Response ID's (like UserID, RunnerID, ...).
- Return `*gitlab.Response` is the API Call response itself
- Return `error` is the API Call response error itself

Example:

```go
func (h MyHandler) ApiCallFunc(page int) ([]int, *gitlab.Response, error) {
	opt := gitlab.UserListOptions{
        Page: page,
    }
    values, resp, err := gitlab.Users.ListUsers(&opt)

    // transfer result value into slice
	var xIds []int
	for _, value := range values {
		xIds = append(xIds, value.ID)
	}
	return xIds, resp, err
}
```
