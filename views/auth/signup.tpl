{{define "head"}}
<style></style>
{{ end }}

{{define "content"}}
<main class="main-content mt-0">
    <section>
        <div
            class="page-header align-items-start min-vh-50 pt-5 pb-11 m-3 border-radius-lg"
            style="background-image: url('/public/img/background.png')"
        >
            <span class="mask bg-gradient-dark opacity-6"></span>
            <div class="container">
                <div class="row justify-content-center">
                    <div class="col-lg-5 text-center mx-auto">
                        <h1 class="text-white mb-2 mt-5">Parrot Disco LTE</h1>
                        <p class="text-lead text-white">
                            Create free account and deploy ready to fly Parrot Disco LTE Dashboard within a few minutes.
                        </p>
                    </div>
                </div>
            </div>
        </div>
        <div class="container">
            <div class="row mt-lg-n10 mt-md-n11 mt-n10">
                <div class="col-xl-4 col-lg-5 col-md-7 mx-auto">
                    <div class="card z-index-0">
                        <div class="card-header text-center pt-4">
                            <h5>Register for free</h5>
                        </div>
                        <div class="card-body">
                            {{if .Alert.Dangers}}
                            {{ range.Alert.Dangers }}
                            <div class="alert alert-danger text-white text-center" role="alert">{{ . }}</div>
                            {{ end }}
                            {{ end }}

                            {{if .Alert.Successes}}
                            {{ range.Alert.Successes }}
                            <div class="alert alert-success text-white text-center" role="alert">{{ . }}</div>
                            {{ end }}
                            {{ end }}

                            <form role="form text-left" method="POST">
                                <div class="mb-3">
                                    <input
                                        type="text"
                                        class="form-control"
                                        placeholder="Name"
                                        aria-label="Name"
                                        aria-describedby="email-addon"
                                        name="name"
                                        required=""
                                    />
                                </div>
                                <div class="mb-3">
                                    <input
                                        type="email"
                                        class="form-control"
                                        placeholder="Email"
                                        aria-label="Email"
                                        aria-describedby="email-addon"
                                        name="email"
                                        required=""
                                    />
                                </div>
                                <div class="mb-3">
                                    <input
                                        type="password"
                                        class="form-control"
                                        placeholder="Password"
                                        aria-label="Password"
                                        aria-describedby="password-addon"
                                        name="password"
                                        required=""
                                        minlength="8"
                                        maxlength="50"
                                    />
                                </div>
                                <div class="form-check form-check-info text-left">
                                    <input
                                        class="form-check-input"
                                        type="checkbox"
                                        value=""
                                        id="flexCheckDefault"
                                        required=""
                                    />
                                    <label class="form-check-label" for="flexCheckDefault">
                                        I declare that I will use the software in accordance with the law and I'm aware
                                        that the software is in the testing phase, therefore errors may occur for which
                                        I won't blame the software author.
                                    </label>
                                </div>
                                <div class="text-center">
                                    <button type="submit" class="btn bg-gradient-dark w-100 my-4 mb-2">Sign up</button>
                                </div>
                                <p class="text-sm mt-3 mb-0">
                                    Already have an account?
                                    <a href="/auth/sign-in" class="text-dark font-weight-bolder">Sign in</a>
                                </p>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
    <footer class="footer mt-4">
        <div class="container">
            <div class="row">
                <div class="col-lg-8 mx-auto text-center mb-4">
                    <a
                        href="https://www.facebook.com/groups/173527569745090"
                        target="_blank"
                        class="text-secondary me-xl-4 me-4"
                    >
                        <span class="text-lg fab fa-facebook"></span>
                    </a>
                    <a href="https://uavpal.slack.com/" target="_blank" class="text-secondary me-xl-4 me-4">
                        <span class="text-lg fab fa-slack"></span>
                    </a>
                    <a
                        href="https://www.youtube.com/@ClassAxion/videos"
                        target="_blank"
                        class="text-secondary me-xl-4 me-4"
                    >
                        <span class="text-lg fab fa-youtube"></span>
                    </a>
                    <a href="https://github.com/ClassAxion" target="_blank" class="text-secondary me-xl-4 me-4">
                        <span class="text-lg fab fa-github"></span>
                    </a>
                </div>
            </div>
            <div class="row">
                <div class="col-8 mx-auto text-center mt-1">
                    <p class="mb-0 text-secondary">
                        Copyright ©
                        <script>
                            document.write(new Date().getFullYear());
                        </script>
                        <br />
                        Made with <i class="fa fa-heart"></i> in Poland
                    </p>
                </div>
            </div>
        </div>
    </footer>
</main>
{{ end }}

{{define "footer"}}
<script></script>
{{ end }}
