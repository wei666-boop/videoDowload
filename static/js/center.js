async function loadRecords(){
    const response=await fetch("/center/record")
    const records=await response.json()

    const tbody=document.getElementById("record-table-body")
    tbody.innerText=""

    records.forEach(record=>{
        const tr=document.createElement("tr")
        tr.innerHTML=`
        <td>${record.id}</td>
        <td>${record.url}</td>
        <td>${record.time}</td>
            `
        tbody.appendChild(tr)
    })
    
}

function clearRecords(){
    const tbody=document.getElementById("record-table-body")
    tbody.innerText=""
    fetch("/center/clear",{
        method:"DELETE"
    })
        .then((res)=>{
            if(!res.ok){
                alert("清空失败")
            }
        })
        .catch(err=>{
            console.error(err)
            alert("清空历史记录时出错: " + err.message)
        })
}