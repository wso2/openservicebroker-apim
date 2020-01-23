package de.evoila.osb.checker

import de.evoila.osb.checker.tests.BindingJUnit5
import de.evoila.osb.checker.tests.CatalogJUnit5
import de.evoila.osb.checker.tests.ProvisionJUnit5
import de.evoila.osb.checker.tests.contract.AuthenticationJUnit5
import de.evoila.osb.checker.tests.contract.ContractJUnit5
import de.evoila.osb.checker.util.TreePrintingListener
import org.apache.commons.cli.DefaultParser
import org.apache.commons.cli.Option
import org.apache.commons.cli.Options
import org.junit.platform.engine.DiscoverySelector
import org.junit.platform.engine.discovery.DiscoverySelectors
import org.junit.platform.launcher.core.LauncherDiscoveryRequestBuilder
import org.junit.platform.launcher.core.LauncherFactory
import org.junit.platform.launcher.listeners.LoggingListener
import org.junit.platform.launcher.listeners.SummaryGeneratingListener
import org.springframework.boot.autoconfigure.SpringBootApplication
import java.io.PrintWriter
import java.util.logging.Level
import kotlin.system.exitProcess

@SpringBootApplication
class Application

fun main(args: Array<String>) {

    val options = Options()
            .apply {
                addOption(
                        Option.builder("cat")
                                .longOpt("catalog")
                                .desc("Indicate if the Catalog Test should run.")
                                .build()
                )
                addOption(
                        Option.builder("prov")
                                .longOpt("provision")
                                .desc("Indicate if the Provision Test should run.")
                                .build()
                )
                addOption(
                        Option.builder("bind")
                                .longOpt("binding")
                                .desc("Indicate if the Binding Test should run.")
                                .build()
                )
                addOption(
                        Option.builder("auth")
                                .longOpt("authentication")
                                .desc("Indicate if the Authentication Test should run.")
                                .build()
                )
                addOption(Option.builder("con")
                        .longOpt("contract")
                        .desc("Indicate if the Contract Test should run.")
                        .build()
                )
            }

    val parser = DefaultParser()
    val commandLine = parser.parse(options, args)

    val selectors = mutableListOf<DiscoverySelector>()

    if (commandLine.hasOption("catalog")) {
        selectors.add(DiscoverySelectors.selectClass(CatalogJUnit5::class.java))
    }

    if (commandLine.hasOption("provision")) {
        selectors.add(DiscoverySelectors.selectClass(ProvisionJUnit5::class.java))
    }

    if (commandLine.hasOption("binding")) {
        selectors.add(DiscoverySelectors.selectClass(BindingJUnit5::class.java))
    }

    if (commandLine.hasOption("authentication")) {
        selectors.add(DiscoverySelectors.selectClass(AuthenticationJUnit5::class.java))
    }

    if (commandLine.hasOption("contract")) {
        selectors.add(DiscoverySelectors.selectClass(ContractJUnit5::class.java))
    }

    if (selectors.isEmpty()) {
        selectors.addAll(
                listOf(DiscoverySelectors.selectClass(CatalogJUnit5::class.java),
                        DiscoverySelectors.selectClass(ProvisionJUnit5::class.java),
                        DiscoverySelectors.selectClass(BindingJUnit5::class.java),
                        DiscoverySelectors.selectClass(AuthenticationJUnit5::class.java),
                        DiscoverySelectors.selectClass(ContractJUnit5::class.java)))
    }

    val summaryGenerator = SummaryGeneratingListener()
    val launcher = LauncherFactory.create()
    val treeLogger = TreePrintingListener()

    launcher.registerTestExecutionListeners(
            LoggingListener.forJavaUtilLogging(Level.INFO),
            treeLogger,
            summaryGenerator
    )

    val request = LauncherDiscoveryRequestBuilder.request()
            .selectors(selectors)
            .build()

    launcher.execute(request)
    summaryGenerator.summary.printTo(PrintWriter(System.out))

    exitProcess(summaryGenerator.summary.failures.size)
}
