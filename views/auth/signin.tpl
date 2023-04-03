{{define "head"}}
<style></style>
{{ end }}

{{define "content"}}
<main class="main-content mt-0">
    <section>
        <div class="page-header min-vh-75">
            <div class="container">
                <div class="row">
                    <div class="col-xl-4 col-lg-5 col-md-6 d-flex flex-column mx-auto">
                        <div class="card card-plain mt-8">
                            <div class="card-header pb-0 text-left bg-transparent">
                                <h3 class="font-weight-bolder text-info text-gradient">Welcome back</h3>
                                <p class="mb-0">Enter your email and password to sign in</p>
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

                                <form role="form" method="POST">
                                    <label>Email</label>
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
                                    <label>Password</label>
                                    <div class="mb-3">
                                        <input
                                            type="password"
                                            class="form-control"
                                            placeholder="Password"
                                            aria-label="Password"
                                            aria-describedby="password-addon"
                                            name="password"
                                            required=""
                                        />
                                    </div>
                                    <div class="text-center">
                                        <button type="submit" class="btn bg-gradient-info w-100 mt-4 mb-0">
                                            Sign in
                                        </button>
                                    </div>
                                </form>
                            </div>
                            <div class="card-footer text-center pt-0 px-lg-2 px-1">
                                <p class="mb-4 text-sm mx-auto">
                                    Don't have an account?
                                    <a href="/auth/sign-up" class="text-info text-gradient font-weight-bold">Sign up</a>
                                </p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-6">
                        <div class="oblique position-absolute top-0 h-100 d-md-block d-none me-n8">
                            <div
                                class="oblique-image bg-cover position-absolute fixed-top ms-auto h-100 z-index-0 ms-n6"
                                style="background-image: url('/public/img/background.png')"
                            ></div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
</main>

<footer class="footer py-5">
    <div class="container">
        <div class="row">
            <div class="col-lg-8 mx-auto text-center mb-4 mt-2">
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
{{ end }}

{{define "footer"}}
<script></script>
{{ end }}
