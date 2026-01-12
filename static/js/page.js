let windowFlag=0

function closeWindow(){
    modal.style.display="none"
}

function controlWindow(){
    const modal=document.getElementById("modal")
    if (windowFlag===1){
        modal.style.display="none"
        windowFlag=0
    }
    else if(windowFlag===0){
        modal.style.display="flex"
        windowFlag=1
    }
}


function videoDownload() {
    const url = document.getElementById("text-url")
    const url64 = btoa(url.value)

    const subtitleDOM = document.getElementById("text-subtitle")

    const thumbnailDOM = document.getElementById("text-thumbnail")

    const formatDOM = document.getElementById("text-format")

    const typeDOM = document.getElementById("select-type")

    const progressDOM=document.getElementById("progress")

    let subtitle="false"
    let thumbnail="false"
    let type="audio"

    let success=true

    if (formatDOM.value === "auto") {
        if (typeDOM.value === "audio") {
            type = "audio"
            subtitle = "false"
            thumbnail = "true"
        } else if (typeDOM.value === "video") {
            type = "video"
            subtitle = "false"
            thumbnail = "true"
        }
    } else {
        if (typeDOM.value === "audio") {
            type = "audio"
            subtitle = "false"
            if (thumbnailDOM.checked) {
                thumbnail= "true"
            } else {
               thumbnail = "false"
            }
        } else if (typeDOM.value === "video") {
            type = "video"
            if (subtitleDOM.checked) {
                subtitle = "true"
            } else {
                subtitle = "false"
            }
            if (thumbnailDOM.checked) {
                thumbnail = "true"
            } else {
                thumbnail = "false"
            }
        }
    }

    progressDOM.innerText="正在配置类型"

    const config = {
        url: url64,
        subtitle: subtitle,
        thumbnail: thumbnail,
        type: type,
    }

    // const baseURL=new URL("http://localhost:5443/dl/api")
    // baseURL.searchParams.append("config",JSON.stringify(config))


    progressDOM.innerText="正在发送请求中"

    url.innerText=""

    fetch("http://localhost:5443/dl/api", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(config)
        //下载文件创建a标签
    }).then(res => {
        if(!res.ok){
            throw new Error(`HTTP Error ${res.status}`)
        }
        progressDOM.innerText="下载中"
        let output="";
        if(res.headers.get("Content-Type")?.includes("video/x-matroska"))
            output="output.mkv"
        else if(res.headers.get("Content-Type")?.includes("video/mp4"))
            output="video.mp4"
        else if(res.headers.get("Content-Type")?.includes("audio/mpeg"))
            output="audio.mp3"
        return res.blob().then(blob=>[blob,output])
    })//返回promise对象
        .then(([blob,output]) => {
            const link = document.createElement("a")
            link.href = URL.createObjectURL(blob)
            link.download = output
            link.click()
            progressDOM.innerText="正在保存文件"
            URL.revokeObjectURL(link.href)
        })
        .catch(error=> {
            console.error(error)
            success=false
        }
        )
        .finally(_=> {
                if (success === false) {
                    progressDOM.innerText = "下载失败"
                } else {
                    progressDOM.innerText = "下载成功"
                }
            }
        )
}

