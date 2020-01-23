package de.evoila.osb.checker.util

import com.google.common.collect.Queues

import org.junit.platform.engine.TestExecutionResult
import org.junit.platform.engine.reporting.ReportEntry
import org.junit.platform.launcher.TestExecutionListener
import org.junit.platform.launcher.TestIdentifier
import org.junit.platform.launcher.TestPlan
import java.io.PrintWriter

class TreePrintingListener : TestExecutionListener {

  private val stack = Queues.newArrayDeque<TreeNode>()

  private val treePrinter = MyTreePrinter(
      out = PrintWriter(System.out),
      theme = Theme.UNICODE
  )

  override fun testPlanExecutionStarted(testPlan: TestPlan) {
    stack.push(TreeNode(testPlan.toString()))
  }

  override fun testPlanExecutionFinished(testPlan: TestPlan) {
    if (stack.isNotEmpty()) {
      treePrinter.print(stack.pop())
    }
  }

  override fun executionStarted(testIdentifier: TestIdentifier) {
    val treeNode = TreeNode(testIdentifier)
    stack.peek()?.addChild(treeNode)
    stack.push(treeNode)
  }

  override fun executionFinished(testIdentifier: TestIdentifier, testExecutionResult: TestExecutionResult) {

    if (stack.isNotEmpty()) {
      stack.pop().setResult(testExecutionResult)
    }
  }

  override fun executionSkipped(testIdentifier: TestIdentifier, reason: String) {
    stack.peek().addChild(TreeNode(testIdentifier, reason))
  }

  override fun reportingEntryPublished(testIdentifier: TestIdentifier, entry: ReportEntry) {
    stack.peek().addReportEntry(entry)
  }

  override fun dynamicTestRegistered(testIdentifier: TestIdentifier) {

  }
}