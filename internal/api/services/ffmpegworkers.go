package services

import (
	"context"
	"encoding/json"
	"io"
	"klms/internal/api/storage/minio"
	"log"
	"os"
	"strings"

	"os/exec"

	sdk "github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)



func Worker() error{

	   conchl,chlerr :=  RabbitConn.Channel()

	   if chlerr != nil {
		   return chlerr
	   }

	   defer conchl.Close()


	   conchl.Qos(
		1,     
		0,
		false,
	)
	

	  msgs,conserr := conchl.Consume(
		"video_queue",
		"worker",
		false,
		false,
		false,
		false,
	    amqp.Table{
			"x-queue-mode": "lazy",
		},
	   )

	   if conserr != nil {
		  return conserr
	   }

	   for msg := range msgs {
             
		    var data map[string]interface{}
		    json.Unmarshal(msg.Body,&data)

			input := data["objectname"].(string)
            foldername := strings.ReplaceAll(data["coursename"].(string), " ", "")
			videoname := data["videoname"].(string)


			object, objfetcherr :=minio.Minio.GetObject(context.Background(),"klms-coursevideos",input,sdk.GetObjectOptions{})

		    if objfetcherr != nil {
				log.Println("cannot fetch from the bucket",objfetcherr)
				continue
			}

             os.MkdirAll("/home/john/Documents/tmp/"+foldername,0777)


				cmd := exec.Command("ffmpeg",
					"-i", "pipe:0",

					// Split video into 2 streams
					"-filter_complex",
					"[0:v]split=2[v0][v1];"+
						"[v0]scale=-2:1080[v0out];"+
						"[v1]scale=-2:360[v1out]",

					// 1080p stream
					"-map", "[v0out]", "-map", "0:a?",
					"-b:v:0", "5000k",

					// 360p stream
					"-map", "[v1out]", "-map", "0:a?",
					"-b:v:1", "800k",

					"-c:v", "libx264",
					"-c:a", "aac",
					"-ac", "2",
					"-ar", "48000",
					"-b:a", "128k",

					// Two video streams (v:0 = 1080p, v:1 = 360p)
					"-var_stream_map", "v:0,a:0 v:1,a:1",

					"-hls_time", "10",
					"-hls_playlist_type", "vod",
					"-master_pl_name", "master.m3u8",

					"-f", "hls",
					"/home/john/Documents/tmp/"+foldername+"/"+"output_%v.m3u8",
				)

				localpath := "/home/john/Documents/tmp/"+foldername


				stdin,stderr := cmd.StdinPipe()

				if stderr != nil {
					log.Println("Cannot get the data from the stdpipe",stderr)
					continue
				}

				cmd.Start()

        
				io.Copy(stdin,object)
				stdin.Close()

                
			    cmd.Wait()

				entries,readerr := os.ReadDir(localpath)

				if readerr != nil {
					log.Println(readerr)
					return readerr
				}

				var objname string

				for _,entry := range entries {

					 name := entry.Name()

					 fullpath := localpath+"/"+name

					 file,openerr := os.Open(fullpath)


					 if openerr != nil {
						log.Println("openning error ",openerr)
						return openerr
					 }

					 fileinfo,infoerr := file.Stat()

					 if infoerr != nil {
						 log.Println("Cannot get the file information",infoerr)
						 file.Close()
						 return infoerr
					 }

					 objname = foldername+"/"+ videoname +"/"+name

					 var contenttype string

					 if strings.HasSuffix(name, ".m3u8") {
						contenttype = "application/vnd.apple.mpegurl"
					} else if strings.HasSuffix(name, ".ts") {
						contenttype = "video/mp2t"
					}

					 _,puterr := minio.Minio.PutObject(context.Background(),"klms-videostreaming",objname,file,fileinfo.Size(),sdk.PutObjectOptions{
						  ContentType: contenttype,
					 })

					 if puterr != nil {
						 log.Println("Cannot put the file into the bucket",puterr)
						 file.Close()
						 return puterr
					 }

					 file.Close()
					 os.Remove(fullpath)

				}

			   os.RemoveAll(localpath)
               msg.Ack(false)


		}

		return nil 
	   
} 

	