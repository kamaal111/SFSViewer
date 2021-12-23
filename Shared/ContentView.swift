//
//  ContentView.swift
//  Shared
//
//  Created by Kamaal M Farah on 22/12/2021.
//

import SwiftUI
import ShrimpExtensions
import SalmonUI

struct ContentView: View {
    @StateObject
    private var sfInfoManager = SFInfoManager()

    private let columns = [
        GridItem(.adaptive(minimum: 80))
    ]

    var body: some View {
        NavigationView {
            Text("Side bar")
            ScrollView {
                LazyVGrid(columns: columns, spacing: 20) {
                    ForEach(sfInfoManager.items, id: \.self) { item in
                        VStack {
                            Image(systemName: item.name)
                                .size(.squared(60))
                            Text(item.name)
                        }
                    }
                }
                .padding(.all)
            }
        }
    }
}

final class SFInfoManager: ObservableObject {

    private let sfInfoHodler = try! SFInfoHodler()

    var items: [SFInfo] {
        sfInfoHodler.items.suffix(100).asArray()
    }

}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView().environment(\.managedObjectContext, PersistenceController.preview.container.viewContext)
    }
}
