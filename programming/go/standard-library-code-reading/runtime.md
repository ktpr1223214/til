---
title: Runtime
---
* hash: a9ea91d571b75d27a49e07f3604b8faa7a48525d
  * tag: go1.14

### G/M/P など構造体

``` go
type schedt struct {
	// accessed atomically. keep at top to ensure alignment on 32-bit systems.
	goidgen   uint64
	lastpoll  uint64 // time of last network poll, 0 if currently polling
	pollUntil uint64 // time to which current poll is sleeping

	lock mutex

	// When increasing nmidle, nmidlelocked, nmsys, or nmfreed, be
	// sure to call checkdead().

	midle        muintptr // idle m's waiting for work
	nmidle       int32    // number of idle m's waiting for work
	nmidlelocked int32    // number of locked m's waiting for work
	mnext        int64    // number of m's that have been created and next M ID
	maxmcount    int32    // maximum number of m's allowed (or die)
	nmsys        int32    // number of system m's not counted for deadlock
	nmfreed      int64    // cumulative number of freed m's

	ngsys uint32 // number of system goroutines; updated atomically

	pidle      puintptr // idle p's
	npidle     uint32
	nmspinning uint32 // See "Worker thread parking/unparking" comment in proc.go.

	// Global runnable queue.
	runq     gQueue
	runqsize int32

	// disable controls selective disabling of the scheduler.
	//
	// Use schedEnableUser to control this.
	//
	// disable is protected by sched.lock.
	disable struct {
		// user disables scheduling of user goroutines.
		user     bool
		runnable gQueue // pending runnable Gs
		n        int32  // length of runnable
	}

	// Global cache of dead G's.
	gFree struct {
		lock    mutex
		stack   gList // Gs with stacks
		noStack gList // Gs without stacks
		n       int32
	}

	// Central cache of sudog structs.
	sudoglock  mutex
	sudogcache *sudog

	// Central pool of available defer structs of different sizes.
	deferlock mutex
	deferpool [5]*_defer

	// freem is the list of m's waiting to be freed when their
	// m.exited is set. Linked through m.freelink.
	freem *m

	gcwaiting  uint32 // gc is waiting to run
	stopwait   int32
	stopnote   note
	sysmonwait uint32
	sysmonnote note

	// safepointFn should be called on each P at the next GC
	// safepoint if p.runSafePointFn is set.
	safePointFn   func(*p)
	safePointWait int32
	safePointNote note

	profilehz int32 // cpu profiling rate

	procresizetime int64 // nanotime() of last change to gomaxprocs
	totaltime      int64 // ∫gomaxprocs dt up to procresizetime
}
```

### main 起動
[proc.go#L112](https://github.com/golang/go/blob/20a838ab94178c55bc4dc23ddc332fce8545a493/src/runtime/proc.go#L112)
``` go
// The main goroutine.
func main() {
	g := getg()

	// Racectx of m0->g0 is used only as the parent of the main goroutine.
	// It must not be used for anything else.
	g.m.g0.racectx = 0

	// Max stack size is 1 GB on 64-bit, 250 MB on 32-bit.
	// Using decimal instead of binary GB and MB because
	// they look nicer in the stack overflow failure message.
	if sys.PtrSize == 8 {
		maxstacksize = 1000000000
	} else {
		maxstacksize = 250000000
	}

	// Allow newproc to start new Ms.
	mainStarted = true

	if GOARCH != "wasm" { // no threads on wasm yet, so no sysmon
		systemstack(func() {
			newm(sysmon, nil)
		})
	}

	// Lock the main goroutine onto this, the main OS thread,
	// during initialization. Most programs won't care, but a few
	// do require certain calls to be made by the main thread.
	// Those can arrange for main.main to run in the main thread
	// by calling runtime.LockOSThread during initialization
    // to preserve the lock.
    // comment: 処理中で以下を呼び出される
    // _g_ := getg()
	// _g_.m.lockedg.set(_g_)
    // _g_.lockedm.set(_g_.m)
    // ここで、_g_ は今の main のはずで、これの実行は OS の main thread で実行されているはず(ここからスタートで)
    // つまり、
    // main goroutine に紐づくOS thread の lockedg を main goroutine に
    // main goroutine の lockedm を main thread に
	lockOSThread()

	if g.m != &m0 {
		throw("runtime.main not on m0")
	}

	doInit(&runtime_inittask) // must be before defer
	if nanotime() == 0 {
		throw("nanotime returning zero")
	}

	// Defer unlock so that runtime.Goexit during init does the unlock too.
	needUnlock := true
	defer func() {
		if needUnlock {
			unlockOSThread()
		}
	}()

	// Record when the world started.
	runtimeInitTime = nanotime()

    // comment: It kicks off the background sweeper goroutine, the background scavenger goroutine, and enables GC.
	gcenable()

	main_init_done = make(chan bool)
	if iscgo {
		...
	}

	doInit(&main_inittask)

	close(main_init_done)

	needUnlock = false
	unlockOSThread()

	if isarchive || islibrary {
		// A program compiled with -buildmode=c-archive or c-shared
		// has a main, but it is not executed.
		return
	}
	fn := main_main // make an indirect call, as the linker doesn't know the address of the main package when laying down the runtime
	fn()
	if raceenabled {
		racefini()
	}

	// Make racy client program work: if panicking on
	// another goroutine at the same time as main returns,
	// let the other goroutine finish printing the panic trace.
	// Once it does, it will exit. See issues 3934 and 20018.
	if atomic.Load(&runningPanicDefers) != 0 {
		// Running deferred functions should not take long.
		for c := 0; c < 1000; c++ {
			if atomic.Load(&runningPanicDefers) == 0 {
				break
			}
			Gosched()
		}
	}
	if atomic.Load(&panicking) != 0 {
		gopark(nil, nil, waitReasonPanicWait, traceEvGoStop, 1)
	}

	exit(0)
	for {
		var x *int32
		*x = 0
	}
}
```
* stack size を 1 GB on 64-bit, 250 MB on 32-bit に設定

### M の生成
proc 生成の os Linux ver
[os_linux.go#L126](https://github.com/golang/go/blob/go1.14/src/runtime/os_linux.go#L126)
``` go
cloneFlags = _CLONE_VM | /* share memory */
		_CLONE_FS | /* share cwd, etc */
		_CLONE_FILES | /* share fd table */
		_CLONE_SIGHAND | /* share sig handler table */
		_CLONE_SYSVSEM | /* share SysV semaphore undo lists (see issue #20763) */
        _CLONE_THREAD /* revisit - okay for now */
...

// comment: 実装は assembly 側かと
//go:noescape
func clone(flags int32, stk, mp, gp, fn unsafe.Pointer) int32

// May run with m.p==nil, so write barriers are not allowed.
//go:nowritebarrier
func newosproc(mp *m) {
	stk := unsafe.Pointer(mp.g0.stack.hi)
	/*
	 * note: strace gets confused if we use CLONE_PTRACE here.
	 */
	if false {
		print("newosproc stk=", stk, " m=", mp, " g=", mp.g0, " clone=", funcPC(clone), " id=", mp.id, " ostk=", &mp, "\n")
	}

	// Disable signals during clone, so that the new thread starts
	// with signals disabled. It will enable them in minit.
	var oset sigset
	sigprocmask(_SIG_SETMASK, &sigset_all, &oset)
	ret := clone(cloneFlags, stk, unsafe.Pointer(mp), unsafe.Pointer(mp.g0), unsafe.Pointer(funcPC(mstart)))
	sigprocmask(_SIG_SETMASK, &oset, nil)

	if ret < 0 {
		print("runtime: failed to create new OS thread (have ", mcount(), " already; errno=", -ret, ")\n")
		if ret == -_EAGAIN {
			println("runtime: may need to increase max user processes (ulimit -u)")
		}
		throw("newosproc")
	}
}
```
* 詳細は兎も角、指定した flag で clone 実行することでプロセス生成(ここでは、スレッド)をする
  * clone は、[sys_linux_amd64.s#L571](https://github.com/golang/go/blob/go1.14/src/runtime/sys_linux_amd64.s#L571)
* こいつは、新しい M 生成に使われる

M 生成
[proc.go#L1699](https://github.com/golang/go/blob/go1.14/src/runtime/proc.go#L1699)
``` go
// Create a new m. It will start off with a call to fn, or else the scheduler.
// fn needs to be static and not a heap allocated closure.
// May run with m.p==nil, so write barriers are not allowed.
//go:nowritebarrierrec
func newm(fn func(), _p_ *p) {
	mp := allocm(_p_, fn)
	mp.nextp.set(_p_)
	mp.sigmask = initSigmask
	if gp := getg(); gp != nil && gp.m != nil && (gp.m.lockedExt != 0 || gp.m.incgo) && GOOS != "plan9" {
        ...
        return
	}
	newm1(mp)
}

// comment: ここで、mp
func newm1(mp *m) {
	if iscgo {
		...
	}
	execLock.rlock() // Prevent process clone.
	newosproc(mp)
	execLock.runlock()
}
```

``` go
// Allocate a new m unassociated with any thread.
// Can use p for allocation context if needed.
// fn is recorded as the new m's m.mstartfn.
//
// This function is allowed to have write barriers even if the caller
// isn't because it borrows _p_.
//
//go:yeswritebarrierrec
func allocm(_p_ *p, fn func()) *m {
    _g_ := getg()
    // comment: getg() の _g_.m.locks++ する
	acquirem() // disable GC because it can be called from sysmon
	if _g_.m.p == 0 {
        // comment: Associate p and the current m. らしいので、
		acquirep(_p_) // temporarily borrow p for mallocs in this function
	}

	// Release the free M list. We need to do this somewhere and
    // this may free up a stack we can use.
    // comment: どっかでやらないとだから、ここでやるってことか
	if sched.freem != nil {
		lock(&sched.lock)
		var newList *m
		for freem := sched.freem; freem != nil; {
            // comment: freeWait == 0 なら安全に消せるらしいので、そうじゃない時は次に進むということかと
			if freem.freeWait != 0 {
				next := freem.freelink
				freem.freelink = newList
				newList = freem
				freem = next
				continue
			}
			stackfree(freem.g0.stack)
			freem = freem.freelink
		}
		sched.freem = newList
		unlock(&sched.lock)
	}

	mp := new(m)
	mp.mstartfn = fn
	mcommoninit(mp)

    // comment: m.g0 に、scheduling 用の goroutine with stack を set
	// In case of cgo or Solaris or illumos or Darwin, pthread_create will make us a stack.
	// Windows and Plan 9 will layout sched stack on OS stack.
	if iscgo || GOOS == "solaris" || GOOS == "illumos" || GOOS == "windows" || GOOS == "plan9" || GOOS == "darwin" {
		mp.g0 = malg(-1)
	} else {
		mp.g0 = malg(8192 * sys.StackGuardMultiplier)
	}
	mp.g0.m = mp

	if _p_ == _g_.m.p.ptr() {
		releasep()
	}
	releasem(_g_.m)

	return mp
}
```
* g0 は特殊な goroutine
  * comment では goroutine with scheduling stack とされる
  * [Go: g0, Special Goroutine](https://medium.com/a-journey-with-go/go-g0-special-goroutine-8c778c6704d8)
  * [Go: What Does a Goroutine Switch Actually Involve?](https://medium.com/a-journey-with-go/go-what-does-a-goroutine-switch-actually-involve-394c202dddb7)

g の allocate(malg)
[proc.go#L3352](https://github.com/golang/go/blob/go1.14/src/runtime/proc.go#L3352)
``` go
// Allocate a new g, with a stack big enough for stacksize bytes.
func malg(stacksize int32) *g {
	newg := new(g)
	if stacksize >= 0 {
		stacksize = round2(_StackSystem + stacksize)
		systemstack(func() {
			newg.stack = stackalloc(uint32(stacksize))
		})
		newg.stackguard0 = newg.stack.lo + _StackGuard
		newg.stackguard1 = ^uintptr(0)
		// Clear the bottom word of the stack. We record g
		// there on gsignal stack during VDSO on ARM and ARM64.
		*(*uintptr)(unsafe.Pointer(newg.stack.lo)) = 0
	}
	return newg
}
```

### sysmon
``` go
// Always runs without a P, so write barriers are not allowed.
//
//go:nowritebarrierrec
func sysmon() {
	lock(&sched.lock)
	sched.nmsys++
	checkdead()
	unlock(&sched.lock)

	lasttrace := int64(0)
	idle := 0 // how many cycles in succession we had not wokeup somebody
	delay := uint32(0)
	for {
		if idle == 0 { // start with 20us sleep...
			delay = 20
		} else if idle > 50 { // start doubling the sleep after 1ms...
			delay *= 2
		}
		if delay > 10*1000 { // up to 10ms
			delay = 10 * 1000
		}
		usleep(delay)
		now := nanotime()
		next, _ := timeSleepUntil()
		if debug.schedtrace <= 0 && (sched.gcwaiting != 0 || atomic.Load(&sched.npidle) == uint32(gomaxprocs)) {
			lock(&sched.lock)
			if atomic.Load(&sched.gcwaiting) != 0 || atomic.Load(&sched.npidle) == uint32(gomaxprocs) {
				if next > now {
					atomic.Store(&sched.sysmonwait, 1)
					unlock(&sched.lock)
					// Make wake-up period small enough
					// for the sampling to be correct.
					sleep := forcegcperiod / 2
					if next-now < sleep {
						sleep = next - now
					}
					shouldRelax := sleep >= osRelaxMinNS
					if shouldRelax {
						osRelax(true)
					}
					notetsleep(&sched.sysmonnote, sleep)
					if shouldRelax {
						osRelax(false)
					}
					now = nanotime()
					next, _ = timeSleepUntil()
					lock(&sched.lock)
					atomic.Store(&sched.sysmonwait, 0)
					noteclear(&sched.sysmonnote)
				}
				idle = 0
				delay = 20
			}
			unlock(&sched.lock)
		}
		// trigger libc interceptors if needed
		if *cgo_yield != nil {
			asmcgocall(*cgo_yield, nil)
		}
		// poll network if not polled for more than 10ms
		lastpoll := int64(atomic.Load64(&sched.lastpoll))
		if netpollinited() && lastpoll != 0 && lastpoll+10*1000*1000 < now {
            atomic.Cas64(&sched.lastpoll, uint64(lastpoll), uint64(now))
            // comment: 関数コメントは、netpoll checks for ready network connections. Returns list of goroutines that become runnable.
			list := netpoll(0) // non-blocking - returns list of goroutines
			if !list.empty() {
				// Need to decrement number of idle locked M's
				// (pretending that one more is running) before injectglist.
				// Otherwise it can lead to the following situation:
				// injectglist grabs all P's but before it starts M's to run the P's,
				// another M returns from syscall, finishes running its G,
				// observes that there is no work to do and no other running M's
				// and reports deadlock.
                incidlelocked(-1)
                // comment: list の中の goroutine を _Grunnable にし、global runqueue に push
                // list は空に
				injectglist(&list)
				incidlelocked(1)
			}
		}
		if next < now {
			// There are timers that should have already run,
			// perhaps because there is an unpreemptible P.
			// Try to start an M to run them.
			startm(nil, false)
		}
		// retake P's blocked in syscalls
		// and preempt long running G's
		if retake(now) != 0 {
			idle = 0
		} else {
			idle++
		}
		// check if we need to force a GC
		if t := (gcTrigger{kind: gcTriggerTime, now: now}); t.test() && atomic.Load(&forcegc.idle) != 0 {
			lock(&forcegc.lock)
			forcegc.idle = 0
			var list gList
			list.push(forcegc.g)
			injectglist(&list)
			unlock(&forcegc.lock)
		}
		if debug.schedtrace > 0 && lasttrace+int64(debug.schedtrace)*1000000 <= now {
			lasttrace = now
			schedtrace(debug.scheddetail > 0)
		}
	}
}
```
### stack 関連

stack 型[runtime2.go#L387](https://github.com/golang/go/blob/go1.14/src/runtime/runtime2.go#L387)
``` go
// Stack describes a Go execution stack.
// The bounds of the stack are exactly [lo, hi),
// with no implicit data structures on either side.
type stack struct {
	lo uintptr
	hi uintptr
}
```

newstack[stack.go#L906](https://github.com/golang/go/blob/go1.14/src/runtime/stack.go#L906)
``` go
// Called from runtime·morestack when more stack is needed.
// Allocate larger stack and relocate to new stack.
// Stack growth is multiplicative, for constant amortized cost.
//
// g->atomicstatus will be Grunning or Gscanrunning upon entry.
// If the scheduler is trying to stop this g, then it will set preemptStop.
//
// This must be nowritebarrierrec because it can be called as part of
// stack growth from other nowritebarrierrec functions, but the
// compiler doesn't check this.
//
//go:nowritebarrierrec
func newstack() {
	thisg := getg()
    ...
	gp := thisg.m.curg
    ...
	sp := gp.sched.sp
	...
	// Allocate a bigger segment and move the stack.
	oldsize := gp.stack.hi - gp.stack.lo
	newsize := oldsize * 2
	if newsize > maxstacksize {
		print("runtime: goroutine stack exceeds ", maxstacksize, "-byte limit\n")
		print("runtime: sp=", hex(sp), " stack=[", hex(gp.stack.lo), ", ", hex(gp.stack.hi), "]\n")
		throw("stack overflow")
	}

	// The goroutine must be executing in order to call newstack,
	// so it must be Grunning (or Gscanrunning).
	casgstatus(gp, _Grunning, _Gcopystack)

	// The concurrent GC will not scan the stack while we are doing the copy since
	// the gp is in a Gcopystack status.
	copystack(gp, newsize)
	if stackDebug >= 1 {
		print("stack grow done\n")
	}
	casgstatus(gp, _Gcopystack, _Grunning)
	gogo(&gp.sched)
}
```
* maxstacksize は main で初期化されている
* stack の copy 中は、_Grunning -> _Gcopystack にしておき、終わったら戻す

* [proc.go#L2444](https://github.com/golang/go/blob/go1.14/src/runtime/proc.go#L2444)
``` go
// One round of scheduler: find a runnable goroutine and execute it.
// Never returns.
func schedule() {
	_g_ := getg()

	if _g_.m.locks != 0 {
		throw("schedule: holding locks")
	}

	if _g_.m.lockedg != 0 {
		stoplockedm()
		execute(_g_.m.lockedg.ptr(), false) // Never returns.
	}

	// We should not schedule away from a g that is executing a cgo call,
	// since the cgo call is using the m's g0 stack.
	if _g_.m.incgo {
		throw("schedule: in cgo")
	}

top:
	pp := _g_.m.p.ptr()
	pp.preempt = false

    // comment:
    // https://github.com/golang/go/blob/20a838ab94178c55bc4dc23ddc332fce8545a493/src/runtime/runtime2.go#L1013
    // sched はここで初期化されている
	if sched.gcwaiting != 0 {
		gcstopm()
		goto top
	}
	if pp.runSafePointFn != 0 {
		runSafePointFn()
	}

	// Sanity check: if we are spinning, the run queue should be empty.
	// Check this before calling checkTimers, as that might call
	// goready to put a ready goroutine on the local run queue.
	if _g_.m.spinning && (pp.runnext != 0 || pp.runqhead != pp.runqtail) {
		throw("schedule: spinning with local work")
	}

	checkTimers(pp, 0)

	var gp *g
	var inheritTime bool

	// Normal goroutines will check for need to wakeP in ready,
	// but GCworkers and tracereaders will not, so the check must
	// be done here instead.
    tryWakeP := false
    // comment: こっからは、どっかで gp が != nil になるまでいくつかのパターンを見に行く形
	if trace.enabled || trace.shutdown {
		gp = traceReader()
		if gp != nil {
			casgstatus(gp, _Gwaiting, _Grunnable)
			traceGoUnpark(gp, 0)
			tryWakeP = true
		}
	}
	if gp == nil && gcBlackenEnabled != 0 {
		gp = gcController.findRunnableGCWorker(_g_.m.p.ptr())
		tryWakeP = tryWakeP || gp != nil
	}
	if gp == nil {
		// Check the global runnable queue once in a while to ensure fairness.
		// Otherwise two goroutines can completely occupy the local runqueue
        // by constantly respawning each other.
        // comment: sched.runq は Global runnable queue. sched.runqsize はそのサイズ
		if _g_.m.p.ptr().schedtick%61 == 0 && sched.runqsize > 0 {
			lock(&sched.lock)
			gp = globrunqget(_g_.m.p.ptr(), 1)
			unlock(&sched.lock)
		}
	}
	if gp == nil {
		gp, inheritTime = runqget(_g_.m.p.ptr())
		// We can see gp != nil here even if the M is spinning,
		// if checkTimers added a local goroutine via goready.
	}
	if gp == nil {
		gp, inheritTime = findrunnable() // blocks until work is available
	}

	// This thread is going to run a goroutine and is not spinning anymore,
	// so if it was marked as spinning we need to reset it now and potentially
	// start a new spinning M.
	if _g_.m.spinning {
		resetspinning()
	}

	if sched.disable.user && !schedEnabled(gp) {
		// Scheduling of this goroutine is disabled. Put it on
		// the list of pending runnable goroutines for when we
		// re-enable user scheduling and look again.
		lock(&sched.lock)
		if schedEnabled(gp) {
			// Something re-enabled scheduling while we
			// were acquiring the lock.
			unlock(&sched.lock)
		} else {
			sched.disable.runnable.pushBack(gp)
			sched.disable.n++
			unlock(&sched.lock)
			goto top
		}
	}

	// If about to schedule a not-normal goroutine (a GCworker or tracereader),
	// wake a P if there is one.
	if tryWakeP {
		if atomic.Load(&sched.npidle) != 0 && atomic.Load(&sched.nmspinning) == 0 {
			wakep()
		}
	}
	if gp.lockedm != 0 {
		// Hands off own p to the locked m,
		// then blocks waiting for a new p.
		startlockedm(gp)
		goto top
	}

	execute(gp, inheritTime)
}
```
* GC
globrunqget
``` go
// Try get a batch of G's from the global runnable queue.
// Sched must be locked.
func globrunqget(_p_ *p, max int32) *g {
...
}
```
* global queue の最初の一つの g を返し、それ以外は _p_ の local queue に積む

stubs.go
``` go
// getg returns the pointer to the current g.
// The compiler rewrites calls to this function into instructions
// that fetch the g directly (from TLS or from the dedicated register).
func getg() *g
```

