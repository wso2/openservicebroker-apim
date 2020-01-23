package de.evoila.osb.checker.util

import com.google.common.collect.Iterables.getOnlyElement
import de.evoila.osb.checker.util.Color.Companion.SKIPPED
import org.junit.platform.commons.util.StringUtils
import org.junit.platform.engine.TestExecutionResult
import org.junit.platform.engine.TestExecutionResult.Status
import org.junit.platform.engine.reporting.ReportEntry
import java.io.PrintWriter

class MyTreePrinter(
    private val theme: Theme,
    private val out: PrintWriter
) {

  fun print(node: TreeNode) {
    out.println(color(Color.CONTAINER, theme.root()))
    print(node, "", true)
    out.flush()
  }

  private fun print(node: TreeNode, indent: String, continuous: Boolean) {
    var varIndent = indent
    if (node.visible) {
      printVisible(node, varIndent, continuous)
    }
    if (node.children.isEmpty()) {
      return
    }
    if (node.visible) {
      varIndent += if (continuous) theme.vertical() else theme.blank()
    }
    val iterator = node.children.iterator()
    while (iterator.hasNext()) {
      print(iterator.next(), varIndent, iterator.hasNext())
    }
  }

  private fun printVisible(node: TreeNode, indent: String, continuous: Boolean) {
    val bullet = if (continuous) theme.entry() else theme.end()
    val prefix = color(Color.CONTAINER, indent + bullet)
    val tabbed = color(Color.CONTAINER, indent + (if (continuous) theme.vertical() else theme.blank()) + theme.blank())
    val caption = colorCaption(node)
    val duration = color(Color.CONTAINER, node.duration.toString() + " ms")
    var icon = color(SKIPPED, theme.skipped())

    node.result?.let {
      val resultColor = Color.valueOf(it)
      icon = color(resultColor, theme.status(it))

    }

    out.print(prefix)
    out.print(" ")
    out.print(caption)

    if (node.duration > 10000 && node.children.isEmpty()) {
      out.print(" ")
      out.print(duration)
    }

    out.print(" ")
    out.print(icon)
    node.result?.let { result -> printThrowable(tabbed, result) }
    node.reason?.let { reason -> printMessage(SKIPPED, tabbed, reason) }
    node.reports.forEach { e -> printReportEntry(tabbed, e) }
    out.println()
  }

  private fun colorCaption(node: TreeNode): String {
    val caption = node.caption

    node.result?.let {
      val resultColor = Color.valueOf(it)
      if (it.status != Status.SUCCESSFUL) {
        return color(resultColor, caption)
      }
    }

    return node.reason?.let {
      color(SKIPPED, caption)
    } ?: color(Color.valueOf(node.identifier!!), caption)


  }

  private fun printThrowable(indent: String, result: TestExecutionResult) {
    if (!result.throwable.isPresent) {
      return
    }
    val throwable = result.throwable.get()
    val message = throwable.message ?: throwable.toString()

    printMessage(Color.FAILED, indent, message)
  }

  private fun printReportEntry(indent: String, reportEntry: ReportEntry) {
    out.println()
    out.print(indent)
    out.print(reportEntry.timestamp)
    val entries = reportEntry.keyValuePairs.entries
    if (entries.size == 1) {
      printReportEntry(" ", getOnlyElement(entries))
      return
    }
    for (entry in entries) {
      out.println()
      printReportEntry(indent + theme.blank(), entry)
    }
  }

  private fun printReportEntry(indent: String, mapEntry: Map.Entry<String, String>) {
    out.print(indent)
    out.print(color(Color.YELLOW, mapEntry.key))
    out.print(" = `")
    out.print(color(Color.GREEN, mapEntry.value))
    out.print("`")
  }

  private fun printMessage(color: Color, indent: String, message: String) {
    val lines = message.split("\\R".toRegex()).dropLastWhile { it.isEmpty() }.toTypedArray()
    out.print(" ")
    out.print(color(color, lines[0]))
    if (lines.size > 1) {
      for (i in 1 until lines.size) {
        out.println()
        out.print(indent)
        if (StringUtils.isNotBlank(lines[i])) {
          val extra = theme.blank()
          out.print(color(color, extra + lines[i]))
        }
      }
    }
  }

  fun color(color: Color, text: String) = "$color$text${Color.NONE}"
}