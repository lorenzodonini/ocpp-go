# ocpp-go

Open Charge Point Protocol implementation in Go.

The library targets modern charge points and central systems, running OCPP version 1.6+.

Given that SOAP will no longer be supported in future versions of OCPP, only OCPP-J is supported in this library.
There are currently no plans of supporting OCPP-S.

## Roadmap

Planned milestones and features:

- [ ] OCPP-J
- [ ] OCPP 1.6
    - [ ] Core Profile
    - [ ] Firmware Profile
    - [ ] Local Auth List Profile
    - [ ] Reservation Profile
    - [ ] Remote Trigger Profile
    - [ ] Smart Charging Profile
- [ ] OCPP 2.0

**Note: The library is still a WIP, therefore expect APIs to change a lot.** 

## Usage

```
go get github.com/lorenzodonini/ocpp-go
```

