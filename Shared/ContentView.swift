//
//  ContentView.swift
//  Shared
//
//  Created by Kamaal M Farah on 22/12/2021.
//

import SwiftUI

struct ContentView: View {
    var body: some View {
        NavigationView {
            Text("Select an item")
                .onAppear(perform: {
                    do {
                        let sfInfo = try SFInfo()
                        print(sfInfo)
                    } catch {
                        print(error)
                    }
                })
        }
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView().environment(\.managedObjectContext, PersistenceController.preview.container.viewContext)
    }
}
