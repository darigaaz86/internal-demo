# internal-demo

## DTM demo cases
 - demonstrate http saga normal flow/retry flow/ rollback flow
 - demonstrate http saga with barrier

## setup
running with diff terminal
```go run app/main.go```
```go run rm-1/main.go```
```go run rm-2/main.go```

## APIs
POST   /api/v1/Trans            
POST   /api/v1/BarrierTransV1       
POST   /api/v1/BarrierTransV2F       
POST   /api/v1/BarrierTransV2S 