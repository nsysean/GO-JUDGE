# GO-JUDGE
Coded this in 2-3 days without any prior GO knowledge, so don't expect much...
+ Please make issues/pull requests if any bugs found/performance tweaks
+ Possible performance tweak:
  @src/handlers/judges/regular-judge.go
  @src/handlers/judges/interactive-judge.go
  Instead of using the -o flag on GNU time command(/usr/bin/time) to write to a file,
  parse the stderr instead, should save up to 100ms per test case,
  but the current issue is that (probably me being ignorant, but) GO randomly
  inserts random characters like these -> *!#\n to the stderr stream,
  help would be appreciated <3
For those who want to run this locally, please make the following changes:

src/handlers/run/init-modes.go
(Around line 13)
old:
	configFile, err := os.Open("../../tasks/" + strconv.Itoa(task) + "/problem.xml")
new:
	configFile, err := os.Open("./assets/tasks/" + strconv.Itoa(task) + "/problem.xml")

src/handlers/run/init-modes.go
(Around line 18)
old: 
	taskAbsPath, _ := filepath.Abs("../../tasks/" + strconv.Itoa(params.Task))
new:
	taskAbsPath, _ := filepath.Abs("./assets/tasks/" + strconv.Itoa(params.Task))

Build:
cd src && go build

Run:
cd src && ./src