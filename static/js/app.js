function notify(msg, msgType) {
    notie.alert({
        type: msgType,
        text: msg,
    })
}

function notifyModal(title, text, icon, confirmationButtonText) {
    Swal.fire({
        title: title,
        html: text,
        icon: icon,
        confirmButtonText: confirmationButtonText
    })
}


function Prompt() {
    let toast = function (c) {
        const{
            msg = '',
            icon = 'success',
            position = 'top-end',

        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    let error = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    async function custom(c) {
        const {
            msg = "",
            title = "",
            showConfirmButton = true,
            icon = "",
        } = c;

        const { value: result } = await Swal.fire({
            icon: icon,
            showConfirmButton:showConfirmButton,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }

            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
            preConfirm: () => {
                if (c.preConfirm !== undefined) {
                    c.preConfirm();
                }
            }
        })

        if ((result) &&
            (result.dismiss !== Swal.DismissReason.cancel) &&
            (result.value !== "") &&
            (c.callback !==undefined)) {
            c.callback(result);
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

function singleRoomChooseDates(c) {
    const {
        elementID = "",
        roomID = "",
        CSRFToken = "",
    } = c;
    document.getElementById(elementID).addEventListener("click", function () {
        let html = `
        <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
            <div class="form-row">
                <div class="col">
                    <div class="form-row" id="reservation-dates-modal">
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                        </div>
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                        </div>

                    </div>
                </div>
            </div>
        </form>
        `;
        attention.custom({
            title: 'Choose y our dates',
            msg: html,
            willOpen: () => {
                const elem = document.getElementById("reservation-dates-modal");
                const rp = new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: true,
                    minDate: new Date(),
                })
            },
            didOpen: () => {
                document.getElementById("start").removeAttribute("disabled");
                document.getElementById("end").removeAttribute("disabled");
            },
            callback: function (result) {
                console.log("called");

                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token", CSRFToken);
                formData.append("room_id", roomID);

                fetch('/search-availability-json', {
                    method: "post",
                    body: formData,
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                msg: '<p>Room is available!</p>' +
                                    '<p><a href="/book-room?id=' +
                                    data.room_id +
                                    '&s=' +
                                    data.start_date +
                                    '&e=' +
                                    data.end_date +
                                    '"class = "btn btn-primary">' +
                                    'Book now!</a></p>',
                                icon: 'success',
                                showConfirmButton: false,
                            })
                        } else {
                            attention.error({
                                msg: "Sorry, this room is not available on this dates >:("
                            })
                        }
                    })
            }
        });
    })
}