package backend.api

import backend.modules.account.presentation.route.accountRoutes
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Application.configureRouting() {
  routing {
    route("/api/health") { get { call.respondText("OK") } }
    route("/api/v1") { accountRoutes() }
  }
}
