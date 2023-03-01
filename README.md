# GoSubSeeker

<center>
[![forthebadge](https://forthebadge.com/generator/?plabel=Made+with&slabel=GO)](https://forthebadge.com) &nbsp;
[![forthebadge](https://forthebadge.com/generator/?plabel=Build+ with&slabel=❤️)](https://forthebadge.com) &nbsp;
[![forthebadge](https://forthebadge.com/generator/?plabel=Open&slabel=Source) &nbsp;
  </center>

<h3 align="center">
    🐞
    <a href="https://github.com/Aryanstha/Code-Chronicles/issues">Report Bug</a> &nbsp; &nbsp;
    ⚙️
    <a href="https://github.com/Aryanstha/Code-Chronicles/issues">Request Feature</a>
</h3>

GoSubSeeker is a command-line tool written in Go programming language that scans a target domain to identify its subdomains. It is designed to be fast, efficient, and easy to use. With GoSubSeeker, you can quickly enumerate subdomains of a given domain, which can be helpful for security researchers and penetration testers.

## Installation

To install GoSubSeeker, you need to have Go installed on your machine. If you don't have it already, you can download it from the official website: https://golang.org/dl/. Once you have Go installed, follow these steps:

- Clone this repository using Git:
  ```bash
git clone https://github.com/your_username/GoSubSeeker.git
  ```
  
- Change into the directory:
```
cd GoSubSeeker
```

- Build the executable:
```
go build
```

- run build file
```
./GoSubSeeker
```

## Usage
```bash
Usage of ./subdomain-scanner -h
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
