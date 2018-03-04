Message protocol
====

Incoming to the ws server messages are prefixed with `app.`. Outgoing from the ws server messages are prefixed with `ws.`

# Outgoing

## Device

### Device discovered
```json
{
    type: "ws.device:discovered",
    data: {
        name: "MOOV",
        id: "45:fg:56"
    }
}
```

### Device status (connected or disconnected)
```json
{
    type: "ws.device:status",
    data: {
        id: "45:fg:56",
        recognizedAs: "hr"  # ["hr", "csc"]
        status: "connected"  # ["connected", "disconnected"]
    }
}
```

### Device measurement (bpm, rps...)
```json
{
    type: "ws.device:measurement",
    data: {
        id: "45:fg:56",
        recognizedAs: "hr",  # ["hr", "csc"]
        bpm: 86,
        revolutions: 4,
        rev_per_sec: 0.005
    }
}
```

# Incoming

## Device

### Scan for devices
```json
{
    type: "app.bt:scan",
    data: {}
}
```

### Stop scanning
```json
{
    type: "app.bt:scan_stop",
    data: {}
}
```

### Connect device
```json
{
    type: "app.device:connect",
    data: {
        id: "45:fg:56"
    }
}
```

