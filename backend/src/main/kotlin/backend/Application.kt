import backend.api.configureRouting
import backend.di.appModule
import backend.di.infraModule
import backend.modules.account.accountModule
import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.application.*
import io.ktor.server.plugins.contentnegotiation.*
import org.koin.ktor.plugin.Koin
import org.koin.logger.slf4jLogger

fun main(args: Array<String>): Unit = io.ktor.server.netty.EngineMain.main(args)

fun Application.module() {
  install(Koin) {
    slf4jLogger()
    modules(appModule, infraModule, accountModule)
  }

  install(ContentNegotiation) { json() }

  configureRouting()
}
