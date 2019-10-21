const express = require('express')
const app = express()

app.use(require('prerender-node'))
app.use(express.static(__dirname  + '/static'))
app.get('*', (req, res) => {
    res.sendfile(__dirname  + '/static/index.html')
})

const port = process.env.TENCENTCLOUD_SERVER_PORT || 8080
app.listen(port, () => console.log('app listening on',port))