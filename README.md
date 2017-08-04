# zymurgauge

Homebrewing automation system

## ToDo:

- [ ] Prometheus instrumentation
- [ ] Coverage Instrumentation (http://damien.lespiau.name/2017/05/building-and-using-coverage.html?utm_source=golangweekly&utm_medium=email)
- [ ] Better graceful shutdown (https://medium.com/@matryer/make-ctrl-c-cancel-the-context-context-bd006a8ad6ff)
- [ ] Re-assess equality (https://golangbot.com/structs/?utm_source=golangweekly&utm_medium=email)
- [ ] Benchmarking
- [ ] Remove global vars
- [ ] Not closing files!
- [ ] DIY routing (https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html)

### UI

- [ ] Stub out Vue.js framework
- [ ] View Beers
- [ ] Add Beer
- [ ] Edit Beer
- [ ] Delete Beer
- [ ] View Fermentations
- [ ] Add Fermentation
- [ ] Edit Fermentation
- [ ] Delete Fermentation
- [ ] View Fermentation History Graph
- [ ] View Chambers
- [ ] Edit Chamber
- [ ] Delete Chamber

### Build

- [ ] Continuous Integration

### Services

#### zymsrv

- [ ] Implement PID controller logic
- [ ] Implement flag package from standard library
- [ ] Dockerize
- [ ] Research User Management / Authentication

#### fermmon

- [ ] Dockerize 
- [ ] Implement flag package from standard library
- [ ] Produce server-side events for herms-sim