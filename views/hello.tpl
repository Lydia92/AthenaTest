<header class="hero-unit">
    <div class="container">
        <div class="row">
            <div class="hero-text">
                <h1>Welcome to the  Beego App!</h1>
                <h2>This is My Test Version</h2>
                <p>node_addr:{{.node_addr}}</p><p>node_port:{{.node_port}}</p>
                      {{range $k,$v := .user}}
                      <tr>
                          <td>{{$v.NodeAddr}}</td>
                          <td>{{$v.NodePort}}</td>
                          <td>{{$v.Ownership}}</td>
                      </tr>
                      {{end}}
            </div>
        </div>
    </div>
</header>