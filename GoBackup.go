/* GoBackup by IMcPwn.
 * Copyright 2016 IMcPwn 

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.

 * For the latest code and contact information visit: http://imcpwn.com
 */
 
// This is basically a replacement to pyWinBackup. I made this because I prefer Go over Python for end-user applications.
// All that's left to be ported over from pyWinBackup is searching and copying files with a specific extension in a specific location (or home directory)

package main

import (
  "fmt"
  "os"
  "os/exec"
  "time"
  "io"
  "flag"
  "strings"
)

func copy_folder(source string, dest string) (err error) {

  sourceinfo, err := os.Stat(source)
  if err != nil {
    return err
  }

  err = os.MkdirAll(dest, sourceinfo.Mode())
  if err != nil {
    return err
  }

  directory, _ := os.Open(source)

  objects, err := directory.Readdir(-1)

  for _, obj := range objects {
    sourcefilepointer := source + "/" + obj.Name()
    destinationfilepointer := dest + "/" + obj.Name()
	
    if obj.IsDir() {
      err = copy_folder(sourcefilepointer, destinationfilepointer)
      if err != nil {
        fmt.Println(err)
      }
    } else {
      err = copy_file(sourcefilepointer, destinationfilepointer)
      if err != nil {
        fmt.Println(err)
      }
    }
  }
  return
}

func copy_file(source string, dest string) (err error) {
  sourcefile, err := os.Open(source)
  if err != nil {
    return err
  }
  defer sourcefile.Close()
  destfile, err := os.Create(dest)
  
  if err != nil {
    return err
  }
  defer destfile.Close()
  
  _, err = io.Copy(destfile, sourcefile)
  if err == nil {
    sourceinfo, err := os.Stat(source)
    if err != nil {
      err = os.Chmod(dest, sourceinfo.Mode())
    }
  }
  return
}

func disconnect(dest string) {
  fmt.Println("[*] Disconnecting from" + dest)
  cmd := exec.Command("net.exe","use",dest,"/delete","/y")
  _, err := cmd.Output()
  if err == nil {
  fmt.Println("[*] Disconnected successfully.")
  } else {
    fmt.Println("[*] Not connected so could not disconnect.")
  }
}

func connect(dest string, user string, pass string, credsReq bool) {
  var cmd *exec.Cmd
  if (credsReq) {
    cmd = exec.Command("net.exe","use",dest,"/user:" + user, pass)
  } else {
    cmd = exec.Command("net.exe","use",dest)
  }
  _, err := cmd.Output()
  if err == nil {
  fmt.Println("[*] Connected to " + dest + " successfully.")
  } else {
    fmt.Println("[X] Error connecting!")
    os.Exit(1)
  }
}


func statusFile(location string, finished bool) {
  file := location + "\\" + "backup-status.txt"
  
  // Check for destination existing
  if _, err := os.Stat(file); os.IsNotExist(err) {
    _, err := os.Create(file)
    if err != nil {
        fmt.Println("[X] " + file + " does not exist and could not be created.")
	    return
    }
  }
  msg := ""
  if (finished) {
    msg = "\r\nBackup completed on " + time.Now().Format("2006-01-02 15:50:13")
  } else {
    msg = "\r\nBackup started on " + time.Now().Format("2006-01-02 15:50:13")
  }
  	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
      fmt.Println("[X] Error writing status file.")
    }

    defer f.Close()

    if _, err = f.WriteString(msg); err != nil {
      fmt.Println("[X] Error writing status file.")
    }
}

func main() {
  destination := ""
  shareUser := ""
  sharePass := ""
  credsRequired := true
  destinationIsShare := true
  userDir := os.Getenv("USERPROFILE")
  user := os.Getenv("USERNAME")
  userDirs := []string{"Documents", "Favorites", "My Documents"}
  otherDirs := []string{userDir + "\\AppData\\Roaming\\Microsoft\\Outlook"}
  destFoldName := user + "-" + time.Now().Format("2006-01-02")
  fullDestination := destination + "\\" + destFoldName
  
  skipUserDirs := false
  skipOtherDirs := false
  
  destFlag := flag.String("dest", "", "The destination of where to back up.")
  shareUserFlag := flag.String("shareUser", "", "The username for the Windows share(ignore this if not backing up to a windows share).")
  sharePassFlag := flag.String("sharePass", "", "The password for the Windows share (ignore this if not backing up to a windows share).")
  credsReqFlag := flag.Bool("credsReq", false, "Specify if credentials are required for the Windows share (ignore this if not backing up to a windows share).")
  destIsShareFlag := flag.Bool("destIsShare", false, "Specify if the destination is a windows share.")
  userFlag := flag.String("user", "", "The username of the account to back up (if blank, back up the current user).")
  userDirsFlag := flag.String("userDirs", "", "The directories inside of the user's home directory to back up (comma separated).")
  otherDirsFlag := flag.String("otherDirs", "", "The full path to any other directories to back up (comma separated).")
  defaultFlag := flag.Bool("default", false, "Use compile time options (if selected, all other options will be ignored).")
  
  flag.Parse()
  
  if !*defaultFlag {
    destination = *destFlag
    shareUser = *shareUserFlag
    sharePass = *sharePassFlag
    credsRequired = *credsReqFlag
    destinationIsShare = *destIsShareFlag
    user = *userFlag
	if *userDirsFlag == "" {
	  skipUserDirs = true
	} else {
	  userDirs = strings.Split(*userDirsFlag, ",")
	}
	if *otherDirsFlag == "" {
	  skipOtherDirs = true
	} else {
	  otherDirs = strings.Split(*otherDirsFlag, ",")
	}
  }
  
  if *userFlag != "USERNAME" {
    userDir = "C:\\Users\\" + user
  }
  
  if (*destFlag == "" && !*defaultFlag) {
    fmt.Println("---- GoBackup by IMcPwn ----\nFor help/the latest version visit http://imcpwn.com\n")
    flag.Usage()
	os.Exit(1)
  }
  
  fmt.Println("---- Welcome to GoBackup by IMcPwn ----\n\nFor help/the latest version visit http://imcpwn.com\n\n")
  
  if (destinationIsShare) {
    disconnect(destination)
    connect(destination, shareUser, sharePass, credsRequired)
  }
  
  // Check for destination existing
  if _, err := os.Stat(fullDestination); os.IsNotExist(err) {
    err := os.Mkdir(fullDestination, 0755)
    if err != nil {
      fmt.Println("[X] " + fullDestination + " does not exist and could not be created.")
	  os.Exit(1)
    }
  }
  
  // Write begin status file
  statusFile(fullDestination, false)
  
  if (!skipUserDirs){ 
    // Copy the userDirs
    for _,currDir := range userDirs {
      fmt.Println("[*] Copying " + userDir + "\\" + currDir)
      err := copy_folder(userDir + "\\" + currDir, fullDestination + "\\" + currDir)
	  if err != nil {
	    fmt.Println("[*] Error copying " + userDir + "\\" + currDir + ". Skipping it.")
	  } else {
	    fmt.Println("[*] Copied " + userDir + "\\" + currDir + " successfully.")
	  }
    }
  } else {
    fmt.Println("[*] No user directories were selected. Not copying any.")
  }
  if (!skipOtherDirs) {
    // Copy the otherDirs
    for _,currDir := range otherDirs {
      fmt.Println("[*] Copying " + currDir)
	  // TODO: Update method of getting the last index of currDir (the actual folder instead of the full path)
      err := copy_folder(currDir, fullDestination + "\\" + strings.Split(currDir, "\\")[len(strings.Split(currDir, "\\")) - 1])
      if err != nil {
	    fmt.Println("[*] Error copying " + currDir + ". Skipping it.")
	  } else {
	    fmt.Println("[*] Copied " + currDir + " successfully.")
	  }
    }
  } else {
    fmt.Println("[*] No other directories were selected. Not copying any.")
  }
  
  // Write end status file
  statusFile(fullDestination, true)
  
  if (destinationIsShare) {
    disconnect(destination)
  }
  //time.Sleep(time.Millisecond * 1300)
}