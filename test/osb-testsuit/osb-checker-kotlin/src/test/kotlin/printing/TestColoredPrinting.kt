package printing

import de.evoila.osb.checker.util.Color
import de.evoila.osb.checker.util.MyTreePrinter
import de.evoila.osb.checker.util.Theme
import org.junit.jupiter.api.Test
import java.io.PrintWriter

class TestColoredPrinting {


  @Test
  fun runColor() {
    val myTreePrinter = MyTreePrinter(
        out = PrintWriter(System.out),
        theme = Theme.ASCII
    )
    println(myTreePrinter.color(Color.RED, "Testing Text"))
  }

}