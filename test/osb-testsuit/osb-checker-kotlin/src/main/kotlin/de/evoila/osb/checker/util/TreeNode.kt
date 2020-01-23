package de.evoila.osb.checker.util

import org.junit.platform.engine.TestExecutionResult
import org.junit.platform.engine.reporting.ReportEntry
import org.junit.platform.launcher.TestIdentifier
import java.util.*

class TreeNode(val caption: String,
               val creation: Long = System.currentTimeMillis(),
               var duration: Long = 0,
               val reason: String? = null,
               val identifier: TestIdentifier? = null,
               var result: TestExecutionResult? = null,
               var reports: MutableList<ReportEntry> = mutableListOf(),
               var children: MutableList<TreeNode> = mutableListOf(),
               var visible: Boolean
) {

  constructor(caption: String) : this(
      caption = caption,
      visible = false
  )


  constructor(testIdentifier: TestIdentifier) : this(
      caption = testIdentifier.displayName,
      identifier = testIdentifier,
      visible = true
  )

  constructor(testIdentifier: TestIdentifier, reason: String) : this(
      caption = testIdentifier.displayName,
      identifier = testIdentifier,
      visible = true,
      reason = reason
  )

  fun addChild(node: TreeNode): TreeNode {
    if (children === Collections.EMPTY_LIST) {
      children = ArrayList()
    }

    children.add(node)
    return this
  }

  fun addReportEntry(reportEntry: ReportEntry): TreeNode {
    reports.add(reportEntry)
    return this
  }

  fun setResult(result: TestExecutionResult): TreeNode {
    this.result = result
    this.duration = System.currentTimeMillis() - creation
    return this
  }
}