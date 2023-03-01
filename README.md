# GoSubSeeker 🔎

<center>
	
![build-with-love](https://user-images.githubusercontent.com/67673221/222180229-23287ecf-e910-410d-94d8-30356d006302.svg)
![made-with-go](https://user-images.githubusercontent.com/67673221/222180277-0bdd8e7b-7852-43a7-8224-daad7d7947bd.svg)
![open-source](https://user-images.githubusercontent.com/67673221/222180309-821620f6-dbb1-441d-96e6-e3fec74cbb02.svg)	
 </center>
 
 <br/>

<h3 align="center">
    🐞
    <a href="https://github.com/Aryanstha/Code-Chronicles/issues">Report Bug</a> &nbsp; &nbsp;
    ⚙️
    <a href="https://github.com/Aryanstha/Code-Chronicles/issues">Request Feature</a>
</h3>

GoSubSeeker is a command-line tool written in Go programming language that scans a target domain to identify its subdomains. It is designed to be fast, efficient, and easy to use. With GoSubSeeker, you can quickly enumerate subdomains of a given domain, which can be helpful for security researchers and penetration testers.

## Installation

To install GoSubSeeker, you need to have Go installed on your machine. If you don't have it already, you can download it from the official website: https://golang.org/dl/. Once you have Go installed, follow these steps:

1. Clone this repository using Git:

```
git clone https://github.com/your_username/GoSubSeeker.git
```
  
2. Change into the directory:
```
cd GoSubSeeker
```

3. Build the executable:
```
go build
```
4. run build file
```
./GoSubSeeker
```

## Usage
```bash
Usage of ./GoSubSeeker -h
  -axfr
		DNS Zone Transfer Protocol (AXFR) of RFC 5936 (default true)
  -d string
		The target Domain
  -depth int
		Scan sub domain depth. range[>=1] (default 1)
  -dns string
		DNS global server (default "8.8.8.8/8.8.4.4")
  -f string
		File contains new line delimited subs (default "dict/subnames_full.txt")
  -fw
		Force scan with wildcard domain (default true)
  -h	Show this help message and exit
  -l string
		The target Domain in file
  -o string
		Output file to write results to (defaults to ./log/{target}).txt
  -t int
		Num of scan threads (default 200)
   ```
   
   ### Show your support

<a href="https://www.buymeacoffee.com/ajty97921p" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-violet.png" alt="Buy Me A Coffee" height= "60px" width= "217px" ></a>
