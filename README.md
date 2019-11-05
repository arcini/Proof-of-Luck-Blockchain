
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

@inproceedings{Milutinovic2016,
 author = {Milutinovic, Mitar and He, Warren and Wu, Howard and Kanwal, Maxinder},
 title = {Proof of Luck: An Efficient Blockchain Consensus Protocol},
 booktitle = {Proceedings of the 1st Workshop on System Software for Trusted Execution},
 series = {SysTEX '16},
 year = {2016},
 isbn = {978-1-4503-4670-2},
 location = {Trento, Italy},
 pages = {2:1--2:6},
 articleno = {2},
 numpages = {6},
 url = {http://doi.acm.org/10.1145/3007788.3007790},
 doi = {10.1145/3007788.3007790},
 acmid = {3007790},
 publisher = {ACM},
 address = {New York, NY, USA},
 keywords = {Blockchain, Consensus Protocol, Intel SGX, Trusted Execution Environments},
} 