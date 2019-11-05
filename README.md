
### Deployment steps:
- clone this repo
- navigate to this directory and rename the example file `mv example.env .env`
- `go run main.go`
- open a new terminal window and `nc localhost 9000`
- input a token amount to set stake
- input a BPM 
- wait a few seconds to see which of the two terminals won 
- open as many terminal windows as you like and `nc localhost 9000` and watch Proof of Stake in action!

Based on proof-of-luck as researched by:
Milutinovic, M., He, W., Wu, H., & Kanwal, M. (2016). Proof of Luck: An Efficient Blockchain Consensus Protocol. In Proceedings of the 1st Workshop on System Software for Trusted Execution (pp. 2:1–2:6). ACM.

