// Open a Kprobe at the entry point of the kernel function and attach the
// pre-compiled program. Each time the kernel function enters, the program
// will increment the execution counter by 1. The read loop below polls this
// map value once per second.
kp, err := link.Kprobe(fn, objs.KprobeExecve, nil)
if err != nil {
	log.Fatalf("opening kprobe: %s", err)
}
defer kp.Close()

// Read loop reporting the total amount of times the kernel
// function was entered, once per second.
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

log.Println("Waiting for events..")

for range ticker.C {
	var value uint64
	if err := objs.KprobeMap.Lookup(mapKey, &value); err != nil {
		log.Fatalf("reading map: %v", err)
	}
	log.Printf("%s called %d times\n", fn, value)
}
