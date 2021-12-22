//
//  SFSViewerApp.swift
//  Shared
//
//  Created by Kamaal M Farah on 22/12/2021.
//

import SwiftUI

@main
struct SFSViewerApp: App {
    let persistenceController = PersistenceController.shared

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(\.managedObjectContext, persistenceController.container.viewContext)
        }
    }
}
