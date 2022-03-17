import axios from "axios"
import { useEffect, useState } from "react"

const MainPage = () => {
    const [info, setInfo] = useState([]);
    const [updateTime, setUpdateTime] = useState();

    const [settings, setSettings] = useState({
        loss: 5,
        insp: 50,
        outsp: 50,
        ping: 20
    })

    const [authInfo, setAuthInfo] = useState({
        login: "",
        pass: ""
    })



    const getInfo = () => {
        let url = "http://mail.leadactiv.ru/getDataNew.php"
        axios.get(url).then((resp) => {
            let users = []
            let now = new Date()
            resp.data.forEach((user) => {
                users.push({
                    name: user.name.replace('_', " "),
                    insp: parseFloat(Math.round(user.download * 100) / 100),
                    outsp: parseFloat(Math.round(user.upload * 100) / 100),
                    ping: parseFloat(user.ping[user.ping.length - 2].split("_ping=")[1]),
                    time: new Date(parseInt(user.ping[user.ping.length - 2].split("_ping=")[0]) * 1000),
                    pings: user.ping
                })
            })
            setUpdateTime(new Date())
            setInfo(users)
        })
    }

    const getSettings = () => {
        let url = "http://mail.leadactiv.ru/settings.php"
        axios.get(url).then((resp) => {
            setSettings(resp.data)
        })
    }

    const sendSettings = (event) => {
        event.preventDefault()
        let url = "http://mail.leadactiv.ru/settings.php"
        axios.post(url, settings).then(() => {
            modalClose(event)
        })
    }

    const loginRequest = (event) => {
        event.preventDefault()
        let url = "http://mail.leadactiv.ru/auth.php"
        axios.post(url, { ...authInfo, action: "login" }).then(() => {
            document.querySelectorAll(".modal-dialog").forEach((dialog) => {
                dialog.style.display = "none";
            })
            document.querySelector('.modal-dialog.settings').style.display = "block"

        }).catch((error) => {
            console.log(error)
        })
    }

    useEffect(() => {
        setTimeout(() => {
            getInfo()
        }, 300000)
    }, [info])

    useEffect(() => {
        getInfo()
        getSettings()
    }, [])

    if (updateTime == undefined) {
        return <div>Загрузка</div>
    }

    const calcLastPingLoss = (pings) => {
        let lastTime = 0
        let lastPingLosstime = ""
        pings.forEach((ping) => {
            let time = new Date(parseInt(ping.split("_")[0]) * 1000)
            if (lastTime > 0 && (time - lastTime > 15000)) {
                lastPingLosstime = time.getMinutes() < 10 ? time.getHours() + ":0" + time.getMinutes() + ":" + time.getSeconds() : time.getHours() + ":" + time.getMinutes() + ":" + time.getSeconds()
            }
            lastTime = time
        })

        return lastPingLosstime
    }

    const pingsConvert = (pings) => {
        let arr = []
        let lastTime = 0
        pings.forEach((ping) => {
            if (ping.length > 2) {
                let time = new Date(parseInt(ping.split("_")[0]) * 1000)
                if (lastTime > 0 && (time - lastTime > 15000)) {
                    arr.push("<b style='color:red'>Потеря пакетов</b>")
                }
                lastTime = time
                let pingStr = ping.split("_")[1]
                let res = time.getMinutes() < 10 ? time.getHours() + ":0" + time.getMinutes() + ":" + time.getSeconds() : time.getHours() + ":" + time.getMinutes() + ":" + time.getSeconds()
                res += " - " + pingStr
                arr.push(res)
            }
        })
        return arr.join("<br>")
    }

    const showPings = (pings) => {
        document.querySelectorAll(".modal-dialog").forEach((dialog) => {
            dialog.style.display = "none";
        })
        document.querySelector('.modal').classList.add('show')
        document.querySelector('.modal-dialog.pings').style.display = "block"
        document.querySelector('.modal').style.display = "block"
        document.querySelector('body').classList.add('modal-open')
        document.querySelector('.modal-dialog.pings .modal-body').innerHTML = pingsConvert(pings)
    }

    const usersInfo = info.map(({ name, insp, outsp, ping, loss, time, pings }) => {
        let diff = (new Date() - time) / 60000

        return <tr key={name}>
            <td>{diff > 20 ? <span className="badge rounded-pill badge-light-danger me-1">Оффлайн</span> : (loss > settings.loss || insp < settings.insp || outsp < settings.outsp || ping > settings.ping) ? <span className="badge rounded-pill badge-light-warning me-1">Онлайн - Низкая скорость</span> : <span className="badge rounded-pill badge-light-success me-1">Онлайн - Хорошая скорость</span>}</td>
            <td><span className="fw-bold">{name}</span></td>
            <td>{insp} м/с</td>
            <td>{outsp} м/с</td>
            <td>{ping}</td>
            <td>{time.getMinutes() < 10 ? time.getHours() + ":0" + time.getMinutes() : time.getHours() + ":" + time.getMinutes()}</td>
            <td>{calcLastPingLoss(pings)}</td>
            <td style={{ cursor: "pointer" }} onClick={() => showPings(pings)}>→</td>
        </tr>
    })

    const settingsClick = (event) => {
        event.preventDefault()
        document.querySelectorAll(".modal-dialog").forEach((dialog) => {
            dialog.style.display = "none";
        })
        if (document.cookie.indexOf('token') == -1) {
            document.querySelector('.modal').classList.add('show')
            document.querySelector('.modal-dialog.password').style.display = "block"
        } else {
            document.querySelector('.modal').classList.add('show')
            document.querySelector('.modal-dialog.settings').style.display = "block"
        }
        document.querySelector('.modal').style.display = "block"
        document.querySelector('body').classList.add('modal-open')
    }

    const modalClose = (event) => {
        event.preventDefault()
        document.querySelector('.modal').classList.remove('show')
        document.querySelector('.modal').style.display = "none"
        document.querySelector('body').classList.remove('modal-open')
    }

    const inputChange = (event) => {
        console.log(event.target.name)
        setSettings(
            {
                ...settings,
                [event.target.name]: parseInt(event.target.value)
            }
        )
    }
    const inputChangeAuth = (event) => {
        console.log(event.target.name)
        setAuthInfo(
            {
                ...authInfo,
                [event.target.name]: event.target.value
            }
        )
    }
    return <div style={{ padding: 20 }}>
        <div className="row" id="basic-table">
            <div className="col-12">
                <div className="card">
                    <div className="card-header" style={{ justifyContent: "center", flexDirection: "column", alignItems: "start" }}>
                        <div style={{ display: "flex", justifyContent: "center", width: "100%" }}>
                            <div style={{ display: "flex", flexDirection: "column", justifyContent: "center", width: "100%" }}>
                                <h1 className="card-title" style={{ fontSize: "1.5rem", flexGrow: 1 }}>Мониторинг качества<br />интернета у операторов</h1>
                                <span>Последнее обновление страницы: {updateTime.getMinutes() < 10 ? updateTime.getHours() + ":0" + updateTime.getMinutes() : updateTime.getHours() + ":" + updateTime.getMinutes()}</span>
                            </div>
                            <img style={{ width: "40px", alignSelf: "flex-end", cursor: "pointer" }} onClick={settingsClick} src="/img/settings.png" />
                        </div>


                    </div>
                    <div className="table-responsive">
                        <table className="table">
                            <thead>
                                <tr>
                                    <th>Статус</th>
                                    <th>Оператор</th>
                                    <th>Входящая скорость</th>
                                    <th>Исходящая скорость</th>
                                    <th>Пинг</th>
                                    <th>Последняя активность</th>
                                    <th>Потеря пакетов</th>
                                    <th>Посмотреть пинги</th>
                                </tr>
                            </thead>
                            <tbody>
                                {usersInfo}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
        <div className="modal fade text-start" id="default" tabIndex="-1" aria-labelledby="myModalLabel1" style={{ display: "none", backgroundColor: "rgba(0, 0, 0, 0.5)" }} aria-hidden="true">
            <div className="modal-dialog modal-dialog-centered password">
                <div className="modal-content">
                    <div className="modal-header">
                        <h4 className="modal-title" id="myModalLabel1">Введите логин и пароль, чтобы открыть форму</h4>
                        <button type="button" onClick={modalClose} className="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                    </div>
                    <div className="modal-body">
                        <form className="form form-horizontal">
                            <div className="row">
                                <div className="col-12">
                                    <div className="mb-1 row">
                                        <div className="col-sm-3">
                                            <label className="col-form-label">Логин</label>
                                        </div>
                                        <div className="col-sm-9">
                                            <input type="text" id="first-name" className="form-control" name="login" onChange={inputChangeAuth} placeholder="" value={authInfo.login} />
                                        </div>
                                    </div>
                                </div>
                                <div className="col-12">
                                    <div className="mb-1 row">
                                        <div className="col-sm-3">
                                            <label className="col-form-label">Пароль</label>
                                        </div>
                                        <div className="col-sm-9">
                                            <input type="password" id="first-name" className="form-control" name="pass" onChange={inputChangeAuth} placeholder="" value={authInfo.pass} />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </form>
                    </div>
                    <div className="modal-footer">
                        <button type="button" className="btn btn-primary waves-effect waves-float waves-light" onClick={loginRequest} data-bs-dismiss="modal">Войти</button>
                    </div>
                </div>
            </div>
            <div className="modal-dialog modal-dialog-centered settings">
                <div className="modal-content">
                    <div className="modal-header">
                        <h4 className="modal-title" id="myModalLabel1">Настройки</h4>
                        <button type="button" onClick={modalClose} className="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                    </div>
                    <div className="modal-body">
                        <form className="form form-horizontal">
                            <div className="row">
                                <div className="col-12">
                                    <div className="mb-1 row">
                                        <div className="col-sm-3">
                                            <label className="col-form-label">Потеря пакетов</label>
                                        </div>
                                        <div className="col-sm-9">
                                            <input type="text" id="first-name" className="form-control" name="loss" onChange={inputChange} placeholder="" value={settings.loss} />
                                        </div>
                                    </div>
                                </div>
                                <div className="col-12">
                                    <div className="mb-1 row">
                                        <div className="col-sm-3">
                                            <label className="col-form-label">Минимальная входящая скорость</label>
                                        </div>
                                        <div className="col-sm-9">
                                            <input type="text" id="first-name" className="form-control" name="insp" onChange={inputChange} placeholder="" value={settings.insp} />
                                        </div>
                                    </div>
                                </div>
                                <div className="col-12">
                                    <div className="mb-1 row">
                                        <div className="col-sm-3">
                                            <label className="col-form-label" >Минимальная исходящая скорость</label>
                                        </div>
                                        <div className="col-sm-9">
                                            <input type="text" id="first-name" className="form-control" name="outsp" onChange={inputChange} placeholder="" value={settings.outsp} />
                                        </div>
                                    </div>
                                </div>
                                <div className="col-12">
                                    <div className="mb-1 row">
                                        <div className="col-sm-3">
                                            <label className="col-form-label" >Максимальный пинг</label>
                                        </div>
                                        <div className="col-sm-9">
                                            <input type="text" id="first-name" className="form-control" name="ping" onChange={inputChange} placeholder="" value={settings.ping} />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </form>
                    </div>
                    <div className="modal-footer">
                        <button type="button" className="btn btn-primary waves-effect waves-float waves-light" onClick={sendSettings} data-bs-dismiss="modal">Сохранить</button>
                    </div>
                </div>
            </div>
            <div className="modal-dialog modal-dialog-centered pings">
                <div className="modal-content">
                    <div className="modal-header">
                        <h4 className="modal-title" id="myModalLabel1">Пинги</h4>
                        <button type="button" onClick={modalClose} className="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                    </div>
                    <div className="modal-body"></div>
                </div>
            </div>
        </div>
    </div>

}

export {
    MainPage
}