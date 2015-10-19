package sdl

/*
#include "sdl_wrapper.h"
#include "events.h"

#if !SDL_VERSION_ATLEAST(2,0,2)
#define SDL_RENDER_TARGETS_RESET (0x2000)
#endif
*/
import "C"
import "unsafe"
import "reflect"

var (
	eventFilterCache  EventFilter
	eventWatchesCache map[EventWatchHandle]*eventFilterCallbackContext = make(map[EventWatchHandle]*eventFilterCallbackContext)
)

//EventType (https://wiki.libsdl.org/SDL_EventType)
type EventType uint32

const (
	FIRSTEVENT              EventType = C.SDL_FIRSTEVENT
	QUIT                    EventType = C.SDL_QUIT
	APP_TERMINATING         EventType = C.SDL_APP_TERMINATING
	APP_LOWMEMORY           EventType = C.SDL_APP_LOWMEMORY
	APP_WILLENTERBACKGROUND EventType = C.SDL_APP_WILLENTERBACKGROUND
	APP_DIDENTERBACKGROUND  EventType = C.SDL_APP_DIDENTERBACKGROUND
	APP_WILLENTERFOREGROUND EventType = C.SDL_APP_WILLENTERFOREGROUND
	APP_DIDENTERFOREGROUND  EventType = C.SDL_APP_DIDENTERFOREGROUND

	/* Window events */
	WINDOWEVENT EventType = C.SDL_WINDOWEVENT
	SYSWMEVENT  EventType = C.SDL_SYSWMEVENT

	/* Keyboard events */
	KEYDOWN     EventType = C.SDL_KEYDOWN
	KEYUP       EventType = C.SDL_KEYUP
	TEXTEDITING EventType = C.SDL_TEXTEDITING
	TEXTINPUT   EventType = C.SDL_TEXTINPUT

	/* Mouse events */
	MOUSEMOTION     EventType = C.SDL_MOUSEMOTION
	MOUSEBUTTONDOWN EventType = C.SDL_MOUSEBUTTONDOWN
	MOUSEBUTTONUP   EventType = C.SDL_MOUSEBUTTONUP
	MOUSEWHEEL      EventType = C.SDL_MOUSEWHEEL

	/* Joystick events */
	JOYAXISMOTION    EventType = C.SDL_JOYAXISMOTION
	JOYBALLMOTION    EventType = C.SDL_JOYBALLMOTION
	JOYHATMOTION     EventType = C.SDL_JOYHATMOTION
	JOYBUTTONDOWN    EventType = C.SDL_JOYBUTTONDOWN
	JOYBUTTONUP      EventType = C.SDL_JOYBUTTONUP
	JOYDEVICEADDED   EventType = C.SDL_JOYDEVICEADDED
	JOYDEVICEREMOVED EventType = C.SDL_JOYDEVICEREMOVED

	/* Game controller events */
	CONTROLLERAXISMOTION     EventType = C.SDL_CONTROLLERAXISMOTION
	CONTROLLERBUTTONDOWN     EventType = C.SDL_CONTROLLERBUTTONDOWN
	CONTROLLERBUTTONUP       EventType = C.SDL_CONTROLLERBUTTONUP
	CONTROLLERDEVICEADDED    EventType = C.SDL_CONTROLLERDEVICEADDED
	CONTROLLERDEVICEREMOVED  EventType = C.SDL_CONTROLLERDEVICEREMOVED
	CONTROLLERDEVICEREMAPPED EventType = C.SDL_CONTROLLERDEVICEREMAPPED

	/* Touch events */
	FINGERDOWN   EventType = C.SDL_FINGERDOWN
	FINGERUP     EventType = C.SDL_FINGERUP
	FINGERMOTION EventType = C.SDL_FINGERMOTION

	/* Gesture events */
	DOLLARGESTURE EventType = C.SDL_DOLLARGESTURE
	DOLLARRECORD  EventType = C.SDL_DOLLARRECORD
	MULTIGESTURE  EventType = C.SDL_MULTIGESTURE

	/* Clipboard events */
	CLIPBOARDUPDATE EventType = C.SDL_CLIPBOARDUPDATE

	/* Drag and drop events */
	DROPFILE EventType = C.SDL_DROPFILE

	/* Render events */
	RENDER_TARGETS_RESET EventType = C.SDL_RENDER_TARGETS_RESET

	USEREVENT EventType = C.SDL_USEREVENT
	LASTEVENT EventType = C.SDL_LASTEVENT
)

const (
	ADDEVENT  = C.SDL_ADDEVENT
	PEEKEVENT = C.SDL_PEEKEVENT
	GETEVENT  = C.SDL_PEEKEVENT
)

const (
	QUERY   = C.SDL_QUERY
	IGNORE  = C.SDL_IGNORE
	DISABLE = C.SDL_DISABLE
	ENABLE  = C.SDL_ENABLE
)

// Event (https://wiki.libsdl.org/SDL_Event)
type Event interface {
	GetType() EventType // EventType is common to all sdl events and
	// while the prefered name for this function may be just Type()
	// that conflicts with field name 'Type' in the struct.
	GetTimestamp() uint32 // All declared sdl events have a timestamp
	// but only the Type field is specified in the sdl docs for events.
	// However, since the timestamp is common IMHO it is a good idea to
	// expose it here also.
}

type CEvent struct {
	Type EventType
	_    [52]byte // padding
}

// CommonEvent struct contains data common to all sdl event types.
// Attaching the GetType and GetTimestamp functions
// ensures conformance to the Event interface for this as
// well as all structs that anom enclose the CommonEvent struct.
type CommonEvent struct {
	Type      EventType
	Timestamp uint32
}

// Returns the event type
func (ce *CommonEvent) GetType() EventType {
	return ce.Type
}

// Returns the event timestamp
func (ce *CommonEvent) GetTimestamp() uint32 {
	return ce.Timestamp
}

// WindowEvent (https://wiki.libsdl.org/SDL_WindowEvent)
type WindowEvent struct {
	CommonEvent
	WindowID uint32
	Event    uint8
	_        uint8 // padding
	_        uint8 // padding
	_        uint8 // padding
	Data1    int32
	Data2    int32
}
type cWindowEvent C.SDL_WindowEvent

//TODO: Key{Up,Down}Event only differ in 'Type' - a boolean field would be enough to distinguish up/down events
// I think this is what you want.

type KeyDownEvent struct {
	CommonEvent
	WindowID uint32
	State    uint8
	Repeat   uint8
	_        uint8 // padding
	_        uint8 // padding
	Keysym   Keysym
}
type cKeyboardEvent C.SDL_KeyboardEvent

type KeyUpEvent struct {
	KeyDownEvent
}

// TextEditingEvent (https://wiki.libsdl.org/SDL_TextEditingEvent)
type TextEditingEvent struct {
	CommonEvent
	WindowID uint32
	Text     [C.SDL_TEXTINPUTEVENT_TEXT_SIZE]byte
	Start    int32
	Length   int32
}
type cTextEditingEvent C.SDL_TextEditingEvent

// TextInputEvent (https://wiki.libsdl.org/SDL_TextInputEvent)
type TextInputEvent struct {
	CommonEvent
	WindowID uint32
	Text     [C.SDL_TEXTINPUTEVENT_TEXT_SIZE]byte
}
type cTextInputEvent C.SDL_TextInputEvent

// MouseMotionEvent (https://wiki.libsdl.org/SDL_MouseMotionEvent)
type MouseMotionEvent struct {
	CommonEvent
	WindowID uint32
	Which    uint32
	State    uint32
	X        int32
	Y        int32
	XRel     int32
	YRel     int32
}
type cMouseMotionEvent C.SDL_MouseMotionEvent

// MouseButtonEvent (https://wiki.libsdl.org/SDL_MouseButtonEvent)
type MouseButtonEvent struct {
	CommonEvent
	WindowID uint32
	Which    uint32
	Button   uint8
	State    uint8
	_        uint8 // padding
	_        uint8 // padding
	X        int32
	Y        int32
}
type cMouseButtonEvent C.SDL_MouseButtonEvent

// MouseWheelEvent (https://wiki.libsdl.org/SDL_MouseWheelEvent)
type MouseWheelEvent struct {
	CommonEvent
	WindowID uint32
	Which    uint32
	X        int32
	Y        int32
}
type cMouseWheelEvent C.SDL_MouseWheelEvent

// JoyAxisEvent (https://wiki.libsdl.org/SDL_JoyAxisEvent)
type JoyAxisEvent struct {
	CommonEvent
	Which JoystickID
	Axis  uint8
	_     uint8 // padding
	_     uint8 // padding
	_     uint8 // padding
	Value int16
	_     uint16 // padding
}
type cJoyAxisEvent C.SDL_JoyAxisEvent

// JoyBallEvent (https://wiki.libsdl.org/SDL_JoyBallEvent)
type JoyBallEvent struct {
	CommonEvent
	Which JoystickID
	Ball  uint8
	_     uint8 // padding
	_     uint8 // padding
	_     uint8 // padding
	XRel  int16
	YRel  int16
}
type cJoyBallEvent C.SDL_JoyBallEvent

// JoyHatEvent (https://wiki.libsdl.org/SDL_JoyHatEvent)
type JoyHatEvent struct {
	CommonEvent
	Which JoystickID
	Hat   uint8
	Value uint8
	_     uint8 // padding
	_     uint8 // padding
}
type cJoyHatEvent C.SDL_JoyHatEvent

// JoyButtonEvent (https://wiki.libsdl.org/SDL_JoyButtonEvent)
type JoyButtonEvent struct {
	CommonEvent
	Which  JoystickID
	Button uint8
	State  uint8
	_      uint8 // padding
	_      uint8 // padding
}
type cJoyButtonEvent C.SDL_JoyButtonEvent

type JoyDeviceEvent struct {
	CommonEvent
	Which JoystickID
}

type ControllerAxisEvent struct {
	CommonEvent
	Which JoystickID
	Axis  uint8
	_     uint8 // padding
	_     uint8 // padding
	_     uint8
	Value int16
	_     uint16 // padding
}
type cControllerAxisEvent C.SDL_ControllerAxisEvent

type ControllerButtonEvent struct {
	CommonEvent
	Which  JoystickID
	Button uint8
	State  uint8
	_      uint8 // padding
	_      uint8 // padding
}
type cControllerButtonEvent C.SDL_ControllerButtonEvent

type ControllerDeviceEvent struct {
	CommonEvent
	Which JoystickID
}
type cControllerDeviceEvent C.SDL_ControllerDeviceEvent

// TouchFingerEvent (https://wiki.libsdl.org/SDL_TouchFingerEvent)
type TouchFingerEvent struct {
	CommonEvent
	TouchID  TouchID
	FingerID FingerID
	X        float32
	Y        float32
	DX       float32
	DY       float32
	Pressure float32
}
type cTouchFingerEvent C.SDL_TouchFingerEvent

// MultiGestureEvent (https://wiki.libsdl.org/SDL_MultiGestureEvent)
type MultiGestureEvent struct {
	CommonEvent
	TouchId    TouchID
	DTheta     float32
	DDist      float32
	X          float32
	Y          float32
	NumFingers uint16
	_          uint16 // padding
}
type cMultiGestureEvent C.SDL_MultiGestureEvent

// DollarGestureEvent (https://wiki.libsdl.org/SDL_DollarGestureEvent)
type DollarGestureEvent struct {
	CommonEvent
	TouchID    TouchID
	GestureID  GestureID
	NumFingers uint32
	Error      float32
	X          float32
	Y          float32
}
type cDollarGestureEvent C.SDL_DollarGestureEvent

// DropEvent (https://wiki.libsdl.org/SDL_DropEvent)
type DropEvent struct {
	CommonEvent
	file unsafe.Pointer
}
type cDropEvent C.SDL_DropEvent

type RenderEvent struct {
	CommonEvent
}

// QuitEvent (https://wiki.libsdl.org/SDL_QuitEvent)
type QuitEvent struct {
	CommonEvent
}

type OSEvent struct {
	CommonEvent
}

type ClipboardEvent struct {
	CommonEvent
}

// UserEvent (https://wiki.libsdl.org/SDL_UserEvent)
type UserEvent struct {
	CommonEvent
	WindowID uint32
	Code     int32
	Data1    unsafe.Pointer
	Data2    unsafe.Pointer
}
type cUserEvent C.SDL_UserEvent

// SysWMEvent (https://wiki.libsdl.org/SDL_SysWMEvent)
type SysWMEvent struct {
	CommonEvent
	msg unsafe.Pointer
}
type cSysWMEvent C.SDL_SysWMEvent

type EventAction C.SDL_eventaction

type EventFilter interface {
	FilterEvent(e Event) bool
}

type eventFilterFunc func(Event) bool

type eventFilterCallbackContext struct {
	filter EventFilter
}

type EventWatchHandle uintptr

func (action EventAction) c() C.SDL_eventaction {
	return C.SDL_eventaction(action)
}

// PumpEvents (https://wiki.libsdl.org/SDL_PumpEvents)
func PumpEvents() {
	C.SDL_PumpEvents()
}

// PeepEvents (https://wiki.libsdl.org/SDL_PeepEvents)
func PeepEvents(events []Event, action EventAction, minType, maxType uint32) int {
	var _events []CEvent = make([]CEvent, len(events))

	if action == ADDEVENT { // the contents of _events matter if they are to be added
		for i := 0; i < len(events); i++ {
			_events[i] = *cEvent(events[i])
		}
	}

	_pevents := (*C.SDL_Event)(unsafe.Pointer(&_events[0]))
	retVal := int(C.SDL_PeepEvents(_pevents, C.int(len(events)), action.c(), C.Uint32(minType), C.Uint32(maxType)))

	if action != ADDEVENT { // put events into slice, events unchanged if action = ADDEVENT
		for i := 0; i < retVal; i++ {
			events[i] = goEvent(&_events[i])
		}
	}
	return retVal
}

// HasEvent (https://wiki.libsdl.org/SDL_HasEvent)
func HasEvent(type_ uint32) bool {
	return C.SDL_HasEvent(C.Uint32(type_)) != 0
}

// HasEvents (https://wiki.libsdl.org/SDL_HasEvents)
func HasEvents(minType, maxType uint32) bool {
	return C.SDL_HasEvents(C.Uint32(minType), C.Uint32(maxType)) != 0
}

// FlushEvent (https://wiki.libsdl.org/SDL_FlushEvent)
func FlushEvent(type_ uint32) {
	C.SDL_FlushEvent(C.Uint32(type_))
}

// FlushEvents (https://wiki.libsdl.org/SDL_FlushEvents)
func FlushEvents(minType, maxType uint32) {
	C.SDL_FlushEvents(C.Uint32(minType), C.Uint32(maxType))
}

// PollEvent (https://wiki.libsdl.org/SDL_PollEvent)
func PollEvent() Event {
	var cevent C.SDL_Event
	ret := C.SDL_PollEvent(&cevent)
	if ret == 0 {
		return nil
	}
	return goEvent((*CEvent)(unsafe.Pointer(&cevent)))
}

func goEvent(cevent *CEvent) Event {
	switch cevent.Type {
	case WINDOWEVENT:
		return (*WindowEvent)(unsafe.Pointer(cevent))
	case SYSWMEVENT:
		return (*SysWMEvent)(unsafe.Pointer(cevent))
	case KEYDOWN:
		return (*KeyDownEvent)(unsafe.Pointer(cevent))
	case KEYUP:
		return (*KeyUpEvent)(unsafe.Pointer(cevent))
	case TEXTEDITING:
		return (*TextEditingEvent)(unsafe.Pointer(cevent))
	case TEXTINPUT:
		return (*TextInputEvent)(unsafe.Pointer(cevent))
	case MOUSEMOTION:
		return (*MouseMotionEvent)(unsafe.Pointer(cevent))
	case MOUSEBUTTONDOWN, MOUSEBUTTONUP:
		return (*MouseButtonEvent)(unsafe.Pointer(cevent))
	case MOUSEWHEEL:
		return (*MouseWheelEvent)(unsafe.Pointer(cevent))
	case JOYAXISMOTION:
		return (*JoyAxisEvent)(unsafe.Pointer(cevent))
	case JOYBALLMOTION:
		return (*JoyBallEvent)(unsafe.Pointer(cevent))
	case JOYHATMOTION:
		return (*JoyHatEvent)(unsafe.Pointer(cevent))
	case JOYBUTTONDOWN, JOYBUTTONUP:
		return (*JoyButtonEvent)(unsafe.Pointer(cevent))
	case JOYDEVICEADDED, JOYDEVICEREMOVED:
		return (*JoyDeviceEvent)(unsafe.Pointer(cevent))
	case CONTROLLERAXISMOTION:
		return (*ControllerAxisEvent)(unsafe.Pointer(cevent))
	case CONTROLLERBUTTONDOWN, CONTROLLERBUTTONUP:
		return (*ControllerButtonEvent)(unsafe.Pointer(cevent))
	case CONTROLLERDEVICEADDED, CONTROLLERDEVICEREMOVED, CONTROLLERDEVICEREMAPPED:
		return (*ControllerDeviceEvent)(unsafe.Pointer(cevent))
	case FINGERDOWN, FINGERUP, FINGERMOTION:
		return (*TouchFingerEvent)(unsafe.Pointer(cevent))
	case DOLLARGESTURE, DOLLARRECORD:
		return (*DollarGestureEvent)(unsafe.Pointer(cevent))
	case MULTIGESTURE:
		return (*MultiGestureEvent)(unsafe.Pointer(cevent))
	case DROPFILE:
		return (*DropEvent)(unsafe.Pointer(cevent))
	case RENDER_TARGETS_RESET:
		return (*RenderEvent)(unsafe.Pointer(cevent))
	case QUIT:
		return (*QuitEvent)(unsafe.Pointer(cevent))
	case USEREVENT:
		return (*UserEvent)(unsafe.Pointer(cevent))
	case CLIPBOARDUPDATE:
		return (*ClipboardEvent)(unsafe.Pointer(cevent))
	default:
		return (*CommonEvent)(unsafe.Pointer(cevent))
	}
}

func cEvent(event Event) *CEvent {
	evv := reflect.ValueOf(event)
	p := evv.Elem()
	return (*CEvent)(unsafe.Pointer(p.UnsafeAddr()))
}

// WaitEventTimeout (https://wiki.libsdl.org/SDL_WaitEventTimeout)
func WaitEventTimeout(timeout int) Event {
	var cevent CEvent
	_event := (*C.SDL_Event)(unsafe.Pointer(&cevent))
	ok := int(C.SDL_WaitEventTimeout(_event, C.int(timeout)))
	if ok == 0 {
		return nil
	}
	return goEvent(&cevent)
}

// WaitEvent (https://wiki.libsdl.org/SDL_WaitEvent)
func WaitEvent() Event {
	var cevent CEvent
	_event := (*C.SDL_Event)(unsafe.Pointer(&cevent))
	ok := int(C.SDL_WaitEvent(_event))
	if ok == 0 {
		return nil
	}
	return goEvent(&cevent)
}

// PushEvent (https://wiki.libsdl.org/SDL_PushEvent)
func PushEvent(event Event) int {
	_event := (*C.SDL_Event)(unsafe.Pointer(cEvent(event)))
	return int(C.SDL_PushEvent(_event))
}

func (ef eventFilterFunc) FilterEvent(e Event) bool {
	return ef(e)
}

func (e *eventFilterCallbackContext) handle() EventWatchHandle {
	return EventWatchHandle(unsafe.Pointer(e))
}

func (e *eventFilterCallbackContext) cptr() unsafe.Pointer {
	return unsafe.Pointer(e)
}

//export goSetEventFilterCallback
func goSetEventFilterCallback(data unsafe.Pointer, e *C.SDL_Event) C.int {
	// No check for eventFilterCache != nil. Why? because it should never be
	// nil since the callback is set/unset based on the last filter being nil
	// /non-nil. If there is an issue, then it should panic here so we can
	// figure out why that is.

	return wrapEventFilterCallback(eventFilterCache, e)
}

//export goEventFilterCallback
func goEventFilterCallback(userdata unsafe.Pointer, e *C.SDL_Event) C.int {
	// same sort of reasoning with goSetEventFilterCallback, userdata should
	// always be non-nil and represent a valid eventFilterCallbackContext. If
	// it doesn't a panic will let us know that there something wrong and the
	// problem can be fixed.

	context := (*eventFilterCallbackContext)(userdata)
	return wrapEventFilterCallback(context.filter, e)
}

func wrapEventFilterCallback(filter EventFilter, e *C.SDL_Event) C.int {
	gev := goEvent((*CEvent)(unsafe.Pointer(e)))
	result := filter.FilterEvent(gev)

	if result {
		return C.SDL_TRUE
	} else {
		return C.SDL_FALSE
	}
}

// SetEventFilter (https://wiki.libsdl.org/SDL_SetEventFilter)
func SetEventFilter(filter EventFilter) {
	if eventFilterCache == nil && filter == nil {
		// nothing to do...
		return
	}

	if eventFilterCache == nil && filter != nil {
		// We had no event filter before and do now; lets set
		// goSetEventFilterCallback() as the event filter.
		C.setEventFilter()
	} else if eventFilterCache != nil && filter == nil {
		// We had an event filter before, but no longer do, lets clear the
		// event filter
		C.clearEventFilter()
	}

	eventFilterCache = filter
}

func SetEventFilterFunc(filterFunc func(Event) bool) {
	SetEventFilter(eventFilterFunc(filterFunc))
}

// GetEventFilter (https://wiki.libsdl.org/SDL_GetEventFilter)
func GetEventFilter() EventFilter {
	return eventFilterCache
}

func isCEventFilterSet() bool {
	return C.SDL_GetEventFilter(nil, nil) == C.SDL_TRUE
}

// FilterEvents (https://wiki.libsdl.org/SDL_FilterEvents)
func FilterEvents(filter EventFilter) {
	context := &eventFilterCallbackContext{filter}
	C.filterEvents(context.cptr())
}

func FilterEventsFunc(filterFunc func(Event) bool) {
	FilterEvents(eventFilterFunc(filterFunc))
}

// AddEventWatch (https://wiki.libsdl.org/SDL_AddEventWatch)
func AddEventWatch(filter EventFilter) EventWatchHandle {
	context := &eventFilterCallbackContext{filter}
	C.addEventWatch(context.cptr())
	eventWatchesCache[context.handle()] = context
	return context.handle()
}

func AddEventWatchFunc(filterFunc func(Event) bool) EventWatchHandle {
	return AddEventWatch(eventFilterFunc(filterFunc))
}

// DelEventWatch (https://wiki.libsdl.org/SDL_DelEventWatch)
func DelEventWatch(handle EventWatchHandle) {
	context, ok := eventWatchesCache[handle]
	if !ok {
		return
	}
	delete(eventWatchesCache, context.handle())
	C.delEventWatch(context.cptr())
}

// EventState (https://wiki.libsdl.org/SDL_EventState)
func EventState(type_ uint32, state int) uint8 {
	return uint8(C.SDL_EventState(C.Uint32(type_), C.int(state)))
}

func GetEventState(type_ uint32) uint8 {
	return uint8(C.SDL_EventState(C.Uint32(type_), QUERY))
}

// RegisterEvents (https://wiki.libsdl.org/SDL_RegisterEvents)
func RegisterEvents(numEvents int) uint32 {
	return uint32(C.SDL_RegisterEvents(C.int(numEvents)))
}
