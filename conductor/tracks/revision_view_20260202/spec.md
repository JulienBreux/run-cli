# Specification - Service Revision View

## Overview
The goal is to provide a detailed view for Cloud Run Service Revisions within the Run CLI TUI. This view will allow users to inspect specific revisions, their traffic allocation, deployment history, and configuration details.

## Requirements
- Display a list of revisions for a selected Service.
- Show traffic allocation percentages for each revision.
- Provide detailed configuration for each revision (image, env vars, resources).
- Display deployment history and status conditions.
- Integrate seamlessly with the existing Service Dashboard.
