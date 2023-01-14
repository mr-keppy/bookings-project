
        let attention = Prompt();



        function Prompt() {
            let toast = function (c) {
                const {
                    msg = "",
                    icon = "success",
                    position = "top-end",
                } = c


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
                    footer = ""
                } = c;

                Swal.fire({
                    title: title,
                    text: msg,
                    icon: "success",
                    footer: footer
                })

            }

            let error = function (c) {

                const {
                    msg = "",
                    title = "",
                    footer = ""
                } = c;

                Swal.fire({
                    title: title,
                    text: msg,
                    icon: "error",
                    footer: footer
                })

            }

            async function custom(c) {
                const {
                    msg = "",
                    title = "",
                    icon = "",
                    showConfirmButton = true
                } = c;


                const { value: result } = await Swal.fire({
                    icon: icon,
                    title: title,
                    html: msg,
                    backdrop: false,
                    showCancelButton: true,
                    focusConfirm: false,
                    showConfirmButton: showConfirmButton,
                    willOpen: () => {
                        if (c.willOpen !== undefined) {
                            c.willOpen();
                        }
                    },
                    preConfirm: () => {
                        return [
                            document.getElementById('start_date').value,
                            document.getElementById('end_date').value
                        ]
                    }
                })

                if (result) {
                    if (result.dismiss !== Swal.DismissReason.cancel) {
                        if (result.value !== "") {
                            if (c.callback !== undefined) {
                                c.callback(result);
                            }
                        }
                        else {
                            c.callback(false);
                        }
                    }
                    else {
                        c.callback(false);
                    }
                }
                //  if (formValues) {
                //      Swal.fire(JSON.stringify(formValues))
                //  }
            }


            return {
                toast: toast,
                success: success,
                error: error,
                custom: custom,
            }

        }
        // Example starter JavaScript for disabling form submissions if there are invalid fields
        (() => {
            'use strict'

            // Fetch all the forms we want to apply custom Bootstrap validation styles to
            const forms = document.querySelectorAll('.needs-validation')

            // Loop over them and prevent submission
            Array.from(forms).forEach(form => {
                form.addEventListener('submit', event => {
                    if (!form.checkValidity()) {
                        event.preventDefault()
                        event.stopPropagation()

                    }

                    form.classList.add('was-validated')
                }, false)
            })
        })()



        function notify(msg, msgType) {
            notie.alert({
                type: msgType,
                text: msg
            })
        }
        function notifyModal(title, msg, msgType) {
            Swal.fire({
                title: title,
                html: msg,
                icon: msgType,
                confirmButtonText: 'Ok'
            })
        }

        { { with .Error} }
        notify("{{.}}", "error")
        { { end } }

        { { with .Flash} }
        notify("{{.}}", "success")
        { { end } }

        { { with .Warning} }
        notify("{{.}}", "warning")
        { { end } }
